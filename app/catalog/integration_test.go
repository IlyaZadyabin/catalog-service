package catalog_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/app/catalog"
	"github.com/IlyaZadyabin/catalog-service/internal/testdb"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogIntegration(t *testing.T) {
	db := testdb.Open(t)
	svc := catalog.NewService(models.NewProductsRepository(db))
	list := catalog.NewCatalogHandler(svc)
	detail := catalog.NewProductHandler(svc)

	t.Run("list filters by category and paginates", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/catalog?category=clothing&limit=2", nil)
		rec := httptest.NewRecorder()
		list.HandleGet(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var response catalog.CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, int64(3), response.Total)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, 2, response.Limit)
		for _, p := range response.Products {
			require.NotNil(t, p.Category)
			assert.Equal(t, "clothing", p.Category.Code)
		}
	})

	t.Run("list filters by price", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=10", nil)
		rec := httptest.NewRecorder()
		list.HandleGet(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var response catalog.CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		require.Greater(t, response.Total, int64(0))
		for _, p := range response.Products {
			assert.Less(t, p.Price, 10.0)
		}
	})

	t.Run("product details inherit missing variant price", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		detail.HandleGetByCode(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var response catalog.ProductDetailsResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "PROD001", response.Code)
		require.NotNil(t, response.Category)
		assert.Equal(t, "clothing", response.Category.Code)
		require.NotEmpty(t, response.Variants)

		product, err := models.NewProductsRepository(db).GetByCode(context.Background(), "PROD001")
		require.NoError(t, err)
		require.Len(t, product.Variants, len(response.Variants))
		for i, v := range product.Variants {
			assert.Equal(t, v.EffectivePrice(product.Price).InexactFloat64(), response.Variants[i].Price)
		}
	})

	t.Run("unknown product is 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/catalog/NOPE", nil)
		req.SetPathValue("code", "NOPE")
		rec := httptest.NewRecorder()
		detail.HandleGetByCode(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
