package models

const (
	DefaultLimit = 10
	MinLimit     = 1
	MaxLimit     = 100
)

type ProductQuery struct {
	Offset        int
	Limit         int
	CategoryCode  string
	PriceLessThan *float64
}

func (q ProductQuery) Normalized() ProductQuery {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < MinLimit {
		q.Limit = MinLimit
	}
	if q.Limit > MaxLimit {
		q.Limit = MaxLimit
	}
	return q
}
