package prices

import (
	"context"
	"time"
)

type PriceDetails struct {
	BrandID   int64
	StartDate time.Time
	EndDate   time.Time
}

type Repository interface {
	Close()
	ListAllForTime(ctx context.Context, t time.Time) ([]PriceDetails, error)
}
