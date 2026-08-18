package catalog

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
)

func parseProductQuery(values url.Values) (models.ProductQuery, error) {
	offset, err := parseOptionalInt(values.Get("offset"), 0)
	if err != nil {
		return models.ProductQuery{}, fmt.Errorf("%w: offset", err)
	}
	limit, err := parseOptionalInt(values.Get("limit"), models.DefaultLimit)
	if err != nil {
		return models.ProductQuery{}, fmt.Errorf("%w: limit", err)
	}

	q := models.ProductQuery{
		Offset:       offset,
		Limit:        limit,
		CategoryCode: values.Get("category"),
	}

	if raw := values.Get("priceLessThan"); raw != "" {
		price, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return models.ProductQuery{}, fmt.Errorf("%w: priceLessThan", apperr.ErrInvalid)
		}
		q.PriceLessThan = &price
	}

	return q.Normalized(), nil
}

func parseOptionalInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, apperr.ErrInvalid
	}
	return parsed, nil
}
