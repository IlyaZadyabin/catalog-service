package categories

import (
	"context"
	"strings"

	"github.com/IlyaZadyabin/catalog-service/app/api"
	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
)

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]api.CategoryDTO, error) {
	categories, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]api.CategoryDTO, len(categories))
	for i, c := range categories {
		out[i] = api.CategoryDTO{Code: c.Code, Name: c.Name}
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, code, name string) (api.CategoryDTO, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	if code == "" || name == "" {
		return api.CategoryDTO{}, apperr.ErrInvalid
	}

	category := &models.Category{Code: code, Name: name}
	if err := s.repo.Create(ctx, category); err != nil {
		return api.CategoryDTO{}, err
	}
	return api.CategoryDTO{Code: category.Code, Name: category.Name}, nil
}
