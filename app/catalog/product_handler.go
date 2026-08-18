package catalog

import (
	"net/http"

	"github.com/IlyaZadyabin/catalog-service/app/api"
)

type ProductDetailsResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category *api.CategoryDTO `json:"category,omitempty"`
	Variants []VariantDTO     `json:"variants,omitempty"`
}

type VariantDTO struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductHandler struct {
	svc *Service
}

func NewProductHandler(svc *Service) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	response, err := h.svc.GetByCode(r.Context(), r.PathValue("code"))
	if err != nil {
		writeCatalogError(w, err, "Failed to fetch product")
		return
	}
	api.OKResponse(w, response)
}
