package prices

import (
	"context"
	"time"
)

type PriceDetails struct {
	BrandID   int64
	StartDate time.Time
	EndDate   time.Time
	PriceList int64
	ProductID int64
	Priority  int64
	Price     int64
	Currency  string
}

type Repository interface {
	Close()
	ListAllForTime(ctx context.Context, brandID, productID int64, t time.Time) ([]PriceDetails, error)
}
