package catalog

import (
	"context"

	"github.com/IlyaZadyabin/catalog-service/models"
)

type ProductRepository interface {
	List(ctx context.Context, q models.ProductQuery) ([]models.Product, int64, error)
	GetByCode(ctx context.Context, code string) (*models.Product, error)
}
