package categories_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/app/categories"
	"github.com/IlyaZadyabin/catalog-service/internal/testdb"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoriesIntegration(t *testing.T) {
	db := testdb.Open(t)
	handler := categories.NewCategoriesHandler(categories.NewService(models.NewCategoryRepository(db)))

	t.Run("lists seeded categories", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleList(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var response categories.CategoriesResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.GreaterOrEqual(t, len(response.Categories), 3)
	})

	t.Run("creates a category", func(t *testing.T) {
		body, _ := json.Marshal(categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		var response api.CategoryDTO
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "electronics", response.Code)
	})

	t.Run("duplicate category is 409", func(t *testing.T) {
		body, _ := json.Marshal(categories.CreateCategoryRequest{Code: "clothing", Name: "Clothing"})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}
