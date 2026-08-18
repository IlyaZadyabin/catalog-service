package catalog

import (
	"errors"
	"net/http"

	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
)

type CatalogResponse struct {
	Products []ProductDTO `json:"products"`
	Total    int64        `json:"total"`
	Offset   int          `json:"offset"`
	Limit    int          `json:"limit"`
}

type ProductDTO struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category *api.CategoryDTO `json:"category,omitempty"`
}

type CatalogHandler struct {
	svc *Service
}

func NewCatalogHandler(svc *Service) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	q, err := parseProductQuery(r.URL.Query())
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid query parameter")
		return
	}

	response, err := h.svc.List(r.Context(), q)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	api.OKResponse(w, response)
}

func writeCatalogError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, apperr.ErrInvalid):
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
	case errors.Is(err, apperr.ErrNotFound):
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
	default:
		api.ErrorResponse(w, http.StatusInternalServerError, fallback)
	}
}
