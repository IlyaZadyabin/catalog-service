package models

import (
	"github.com/shopspring/decimal"
)

// Variant represents a product variant in the catalog.
type Variant struct {
	ID        uint             `gorm:"primaryKey"`
	ProductID uint             `gorm:"not null"`
	Name      string           `gorm:"not null"`
	SKU       string           `gorm:"uniqueIndex;not null"`
	Price     *decimal.Decimal `gorm:"type:decimal(10,2)"`
}

func (v *Variant) TableName() string {
	return "product_variants"
}

func (v *Variant) EffectivePrice(productPrice decimal.Decimal) decimal.Decimal {
	if v.Price == nil {
		return productPrice
	}
	return *v.Price
}
