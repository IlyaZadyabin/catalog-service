package models

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProductQueryNormalized(t *testing.T) {
	assert.Equal(t, 0, ProductQuery{Offset: -3, Limit: 10}.Normalized().Offset)
	assert.Equal(t, MinLimit, ProductQuery{Limit: 0}.Normalized().Limit)
	assert.Equal(t, MaxLimit, ProductQuery{Limit: 500}.Normalized().Limit)
	assert.Equal(t, 7, ProductQuery{Limit: 7}.Normalized().Limit)
}

func TestVariantEffectivePrice(t *testing.T) {
	product := decimal.NewFromFloat(10.99)
	own := decimal.NewFromFloat(11.99)
	zero := decimal.Zero

	withOwn := Variant{Price: &own}
	withZero := Variant{Price: &zero}
	without := Variant{}

	assert.Equal(t, "11.99", withOwn.EffectivePrice(product).StringFixed(2))
	assert.Equal(t, "0.00", withZero.EffectivePrice(product).StringFixed(2))
	assert.Equal(t, "10.99", without.EffectivePrice(product).StringFixed(2))
}
