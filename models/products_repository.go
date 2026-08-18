package models

import (
	"context"

	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{db: db}
}

func (r *ProductsRepository) List(ctx context.Context, q ProductQuery) ([]Product, int64, error) {
	q = q.Normalized()

	var products []Product
	var total int64

	query := r.db.WithContext(ctx).Model(&Product{}).Preload("Category")

	if q.CategoryCode != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", q.CategoryCode)
	}

	if q.PriceLessThan != nil {
		query = query.Where("products.price < ?", *q.PriceLessThan)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("products.id").Offset(q.Offset).Limit(q.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetByCode(ctx context.Context, code string) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("code = ?", code).
		First(&product).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &product, nil
}
