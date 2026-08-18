package categories

import (
	"context"

	"github.com/IlyaZadyabin/catalog-service/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) error
}
