package catalog

import (
	"context"

	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
)

type Service struct {
	products ProductRepository
}

func NewService(products ProductRepository) *Service {
	return &Service{products: products}
}

func (s *Service) List(ctx context.Context, q models.ProductQuery) (CatalogResponse, error) {
	q = q.Normalized()
	products, total, err := s.products.List(ctx, q)
	if err != nil {
		return CatalogResponse{}, err
	}
	return CatalogResponse{
		Products: mapProducts(products),
		Total:    total,
		Offset:   q.Offset,
		Limit:    q.Limit,
	}, nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (ProductDetailsResponse, error) {
	if code == "" {
		return ProductDetailsResponse{}, apperr.ErrInvalid
	}
	product, err := s.products.GetByCode(ctx, code)
	if err != nil {
		return ProductDetailsResponse{}, err
	}
	return mapProductDetails(product), nil
}
