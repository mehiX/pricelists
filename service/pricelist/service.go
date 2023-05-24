package pricelist

import (
	"context"
	"fmt"
	"time"

	"github.com/mehix/pricelists/domain/prices"
	"github.com/mehix/pricelists/domain/prices/h2"
)

type PriceDetails struct {
	StartDate time.Time
	EndDate   time.Time
	Price     int64
	Currency  string
	BrandName string
}

type Service interface {
	Close()
	ProductPriceForTime(ctx context.Context, brandName string, productID int64, t time.Time) (PriceDetails, error)
}

type service struct {
	repo prices.Repository
}

func NewService() Service {
	return &service{}
}

func NewServiceForH2(url string) (Service, error) {
	repo, err := h2.NewRepository(url)
	if err != nil {
		return nil, err
	}

	return &service{repo: repo}, nil
}

func (s *service) Close() {
	if s.repo != nil {
		s.repo.Close()
	}
}

func (s *service) ProductPriceForTime(ctx context.Context, brandName string, productID int64, t time.Time) (PriceDetails, error) {
	if s.repo == nil {
		return PriceDetails{}, fmt.Errorf("no datasource connected")
	}

	return PriceDetails{}, nil
}
