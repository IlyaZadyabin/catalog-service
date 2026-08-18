package catalog

import (
	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/models"
)

func mapProducts(products []models.Product) []ProductDTO {
	dtos := make([]ProductDTO, len(products))
	for i, p := range products {
		dtos[i] = mapProduct(p)
	}
	return dtos
}

func mapProduct(p models.Product) ProductDTO {
	dto := ProductDTO{
		Code:  p.Code,
		Price: p.Price.InexactFloat64(),
	}
	if p.Category != nil {
		dto.Category = mapCategory(p.Category)
	}
	return dto
}

func mapProductDetails(p *models.Product) ProductDetailsResponse {
	dto := ProductDetailsResponse{
		Code:  p.Code,
		Price: p.Price.InexactFloat64(),
	}
	if p.Category != nil {
		dto.Category = mapCategory(p.Category)
	}
	if len(p.Variants) == 0 {
		return dto
	}
	dto.Variants = make([]VariantDTO, len(p.Variants))
	for i, v := range p.Variants {
		dto.Variants[i] = VariantDTO{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.EffectivePrice(p.Price).InexactFloat64(),
		}
	}
	return dto
}

func mapCategory(c *models.Category) *api.CategoryDTO {
	return &api.CategoryDTO{Code: c.Code, Name: c.Name}
}
