package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeProducts struct {
	listFn    func(ctx context.Context, q models.ProductQuery) ([]models.Product, int64, error)
	getFn     func(ctx context.Context, code string) (*models.Product, error)
	lastQuery models.ProductQuery
}

func (f *fakeProducts) List(ctx context.Context, q models.ProductQuery) ([]models.Product, int64, error) {
	f.lastQuery = q
	if f.listFn != nil {
		return f.listFn(ctx, q)
	}
	return nil, 0, nil
}

func (f *fakeProducts) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	if f.getFn != nil {
		return f.getFn(ctx, code)
	}
	return nil, errors.New("unexpected GetByCode")
}

func clothingProduct() models.Product {
	return models.Product{
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: &models.Category{
			Code: "clothing",
			Name: "Clothing",
		},
	}
}

func TestCatalogHandler_HandleGet(t *testing.T) {
	t.Run("successful request with default pagination", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, q models.ProductQuery) ([]models.Product, int64, error) {
				assert.Equal(t, 0, q.Offset)
				assert.Equal(t, 10, q.Limit)
				assert.Empty(t, q.CategoryCode)
				assert.Nil(t, q.PriceLessThan)
				return []models.Product{
					clothingProduct(),
					{
						Code:  "PROD002",
						Price: decimal.NewFromFloat(12.49),
						Category: &models.Category{
							Code: "shoes",
							Name: "Shoes",
						},
					},
				}, 2, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Products, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 0, response.Offset)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, 10.99, response.Products[0].Price)
		require.NotNil(t, response.Products[0].Category)
		assert.Equal(t, "clothing", response.Products[0].Category.Code)
	})

	t.Run("request with custom pagination", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, _ models.ProductQuery) ([]models.Product, int64, error) {
				return []models.Product{}, 50, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?offset=20&limit=5", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, 20, response.Offset)
		assert.Equal(t, 5, response.Limit)
		assert.Equal(t, int64(50), response.Total)
		assert.Equal(t, 20, repo.lastQuery.Offset)
		assert.Equal(t, 5, repo.lastQuery.Limit)
	})

	t.Run("request with category filter", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, q models.ProductQuery) ([]models.Product, int64, error) {
				assert.Equal(t, "clothing", q.CategoryCode)
				return []models.Product{clothingProduct()}, 1, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?category=clothing", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Products, 1)
	})

	t.Run("request with price filter", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, q models.ProductQuery) ([]models.Product, int64, error) {
				require.NotNil(t, q.PriceLessThan)
				assert.Equal(t, 10.0, *q.PriceLessThan)
				return []models.Product{{Code: "PROD003", Price: decimal.NewFromFloat(8.75)}}, 1, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=10", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Products, 1)
	})

	t.Run("invalid priceLessThan", func(t *testing.T) {
		repo := &fakeProducts{}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?priceLessThan=cheap", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Invalid query parameter", response["error"])
	})

	t.Run("limit validation - minimum", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, _ models.ProductQuery) ([]models.Product, int64, error) {
				return []models.Product{}, 0, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?limit=0", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, 1, response.Limit)
		assert.Equal(t, 1, repo.lastQuery.Limit)
	})

	t.Run("limit validation - maximum", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, _ models.ProductQuery) ([]models.Product, int64, error) {
				return []models.Product{}, 0, nil
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog?limit=200", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CatalogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, 100, response.Limit)
		assert.Equal(t, 100, repo.lastQuery.Limit)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &fakeProducts{
			listFn: func(_ context.Context, _ models.ProductQuery) ([]models.Product, int64, error) {
				return nil, 0, errors.New("database error")
			},
		}
		handler := NewCatalogHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Failed to fetch products", response["error"])
	})
}

func TestCatalogHandler_InvalidQuery(t *testing.T) {
	handler := NewCatalogHandler(NewService(&fakeProducts{}))

	for _, raw := range []string{"/catalog?limit=abc", "/catalog?offset=abc"} {
		req := httptest.NewRequest(http.MethodGet, raw, nil)
		rec := httptest.NewRecorder()
		handler.HandleGet(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code, raw)
	}
}
