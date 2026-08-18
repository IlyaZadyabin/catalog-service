package categories

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
)

type CategoriesResponse struct {
	Categories []api.CategoryDTO `json:"categories"`
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesHandler struct {
	svc *Service
}

func NewCategoriesHandler(svc *Service) *CategoriesHandler {
	return &CategoriesHandler{svc: svc}
}

func (h *CategoriesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	api.OKResponse(w, CategoriesResponse{Categories: items})
}

func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	created, err := h.svc.Create(r.Context(), req.Code, req.Name)
	if err != nil {
		writeCategoryError(w, err)
		return
	}
	api.CreatedResponse(w, created)
}

func writeCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperr.ErrInvalid):
		api.ErrorResponse(w, http.StatusBadRequest, "Code and name are required")
	case errors.Is(err, apperr.ErrConflict):
		api.ErrorResponse(w, http.StatusConflict, "Category already exists")
	default:
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
	}
}
