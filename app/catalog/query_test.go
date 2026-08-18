package catalog

import (
	"net/url"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/internal/apperr"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProductQuery(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		q, err := parseProductQuery(url.Values{})
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
		assert.Equal(t, models.DefaultLimit, q.Limit)
		assert.Empty(t, q.CategoryCode)
		assert.Nil(t, q.PriceLessThan)
	})

	t.Run("clamps negative offset", func(t *testing.T) {
		q, err := parseProductQuery(url.Values{"offset": {"-4"}})
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
	})

	t.Run("rejects invalid price", func(t *testing.T) {
		_, err := parseProductQuery(url.Values{"priceLessThan": {"nope"}})
		require.ErrorIs(t, err, apperr.ErrInvalid)
	})

	t.Run("rejects invalid limit", func(t *testing.T) {
		_, err := parseProductQuery(url.Values{"limit": {"abc"}})
		require.ErrorIs(t, err, apperr.ErrInvalid)
	})

	t.Run("rejects invalid offset", func(t *testing.T) {
		_, err := parseProductQuery(url.Values{"offset": {"abc"}})
		require.ErrorIs(t, err, apperr.ErrInvalid)
	})
}
