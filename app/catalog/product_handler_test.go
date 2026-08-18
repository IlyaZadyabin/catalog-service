package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func detailedProduct() *models.Product {
	return &models.Product{
		Code:  "PROD001",
		Price: decimal.NewFromFloat(10.99),
		Category: &models.Category{
			Code: "clothing",
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: decPtr("11.99")},
			{Name: "Variant B", SKU: "SKU001B", Price: nil},
		},
	}
}

func decPtr(s string) *decimal.Decimal {
	d := decimal.RequireFromString(s)
	return &d
}

func TestProductHandler_HandleGetByCode(t *testing.T) {
	t.Run("successful request with category and variants", func(t *testing.T) {
		repo := &fakeProducts{
			getFn: func(_ context.Context, code string) (*models.Product, error) {
				assert.Equal(t, "PROD001", code)
				return detailedProduct(), nil
			},
		}
		handler := NewProductHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response ProductDetailsResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 10.99, response.Price)
		require.NotNil(t, response.Category)
		assert.Equal(t, "clothing", response.Category.Code)
		assert.Equal(t, "Clothing", response.Category.Name)
		require.Len(t, response.Variants, 2)
		assert.Equal(t, 11.99, response.Variants[0].Price)
		assert.Equal(t, 10.99, response.Variants[1].Price)
	})

	t.Run("successful request without category", func(t *testing.T) {
		repo := &fakeProducts{
			getFn: func(_ context.Context, _ string) (*models.Product, error) {
				return &models.Product{
					Code:     "PROD002",
					Price:    decimal.NewFromFloat(12.49),
					Category: nil,
					Variants: []models.Variant{},
				}, nil
			},
		}
		handler := NewProductHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog/PROD002", nil)
		req.SetPathValue("code", "PROD002")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response ProductDetailsResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "PROD002", response.Code)
		assert.Nil(t, response.Category)
		assert.Empty(t, response.Variants)
	})

	t.Run("product not found", func(t *testing.T) {
		repo := &fakeProducts{
			getFn: func(_ context.Context, _ string) (*models.Product, error) {
				return nil, apperr.ErrNotFound
			},
		}
		handler := NewProductHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog/INVALID", nil)
		req.SetPathValue("code", "INVALID")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Product not found", response["error"])
	})

	t.Run("missing product code", func(t *testing.T) {
		called := false
		repo := &fakeProducts{
			getFn: func(_ context.Context, _ string) (*models.Product, error) {
				called = true
				return nil, apperr.ErrNotFound
			},
		}
		handler := NewProductHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog/", nil)
		req.SetPathValue("code", "")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.False(t, called)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Product code is required", response["error"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &fakeProducts{
			getFn: func(_ context.Context, _ string) (*models.Product, error) {
				return nil, assert.AnError
			},
		}
		handler := NewProductHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
		req.SetPathValue("code", "PROD001")
		rec := httptest.NewRecorder()
		handler.HandleGetByCode(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Failed to fetch product", response["error"])
	})
}

func TestMapProductDetails(t *testing.T) {
	dto := mapProductDetails(detailedProduct())
	assert.Equal(t, "PROD001", dto.Code)
	assert.Equal(t, 10.99, dto.Price)
	require.NotNil(t, dto.Category)
	assert.Equal(t, "clothing", dto.Category.Code)
	require.Len(t, dto.Variants, 2)
	assert.Equal(t, 11.99, dto.Variants[0].Price)
	assert.Equal(t, 10.99, dto.Variants[1].Price)
}
