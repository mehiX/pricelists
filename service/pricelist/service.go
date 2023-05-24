package pricelist

import (
	"context"
	"time"
)

type PriceDetails struct {
	StartDate time.Time
	EndDate   time.Time
	Price     int64
	Currency  string
	BrancName string
}

type Service interface {
	ProductPriceForTime(ctx context.Context, brandName string, productID int64, t time.Time) (PriceDetails, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) ProductPriceForTime(ctx context.Context, brandName string, productID int64, t time.Time) (PriceDetails, error) {
	return PriceDetails{}, nil
}
