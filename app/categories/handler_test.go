package categories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCategories struct {
	listFn   func(ctx context.Context) ([]models.Category, error)
	createFn func(ctx context.Context, category *models.Category) error
}

func (f *fakeCategories) List(ctx context.Context) ([]models.Category, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f *fakeCategories) Create(ctx context.Context, category *models.Category) error {
	if f.createFn != nil {
		return f.createFn(ctx, category)
	}
	return nil
}

func TestHandler_HandleList(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		repo := &fakeCategories{
			listFn: func(_ context.Context) ([]models.Category, error) {
				return []models.Category{
					{Code: "clothing", Name: "Clothing"},
					{Code: "shoes", Name: "Shoes"},
					{Code: "accessories", Name: "Accessories"},
				}, nil
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleList(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CategoriesResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Categories, 3)
		assert.Equal(t, "clothing", response.Categories[0].Code)
		assert.Equal(t, "Clothing", response.Categories[0].Name)
	})

	t.Run("empty categories list", func(t *testing.T) {
		repo := &fakeCategories{
			listFn: func(_ context.Context) ([]models.Category, error) {
				return []models.Category{}, nil
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleList(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response CategoriesResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Empty(t, response.Categories)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &fakeCategories{
			listFn: func(_ context.Context) ([]models.Category, error) {
				return nil, errors.New("database error")
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		handler.HandleList(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Failed to fetch categories", response["error"])
	})
}

func TestHandler_HandleCreate(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		repo := &fakeCategories{
			createFn: func(_ context.Context, cat *models.Category) error {
				assert.Equal(t, "electronics", cat.Code)
				assert.Equal(t, "Electronics", cat.Name)
				return nil
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		body, _ := json.Marshal(CreateCategoryRequest{Code: "electronics", Name: "Electronics"})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var response api.CategoryDTO
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "electronics", response.Code)
		assert.Equal(t, "Electronics", response.Name)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		called := false
		repo := &fakeCategories{
			createFn: func(_ context.Context, _ *models.Category) error {
				called = true
				return nil
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("invalid json")))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.False(t, called)
	})

	t.Run("missing required fields", func(t *testing.T) {
		repo := &fakeCategories{}
		handler := NewCategoriesHandler(NewService(repo))

		for _, tc := range []CreateCategoryRequest{
			{Code: "", Name: "Electronics"},
			{Code: "electronics", Name: ""},
			{Code: "", Name: ""},
		} {
			body, _ := json.Marshal(tc)
			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			handler.HandleCreate(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			var response map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			assert.Equal(t, "Code and name are required", response["error"])
		}
	})

	t.Run("duplicate category", func(t *testing.T) {
		repo := &fakeCategories{
			createFn: func(_ context.Context, _ *models.Category) error {
				return apperr.ErrConflict
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		body, _ := json.Marshal(CreateCategoryRequest{Code: "clothing", Name: "Clothing"})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		var response map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, "Category already exists", response["error"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &fakeCategories{
			createFn: func(_ context.Context, _ *models.Category) error {
				return errors.New("database error")
			},
		}
		handler := NewCategoriesHandler(NewService(repo))

		body, _ := json.Marshal(CreateCategoryRequest{Code: "electronics", Name: "Electronics"})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.HandleCreate(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
