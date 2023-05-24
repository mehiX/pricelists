package pricelist

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mehix/pricelists/domain/brands"
	h2Brands "github.com/mehix/pricelists/domain/brands/h2"
	"github.com/mehix/pricelists/domain/prices"
	h2Prices "github.com/mehix/pricelists/domain/prices/h2"
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
	pricesRepo prices.Repository
	brandsRepo brands.Repository
}

func NewService() Service {
	return &service{}
}

// NewServiceForH2 returns a service that holds a connection to an H2 instance.
// For now we assume that all repositories point to the same database instance.
// Pointing different interfaces to different instances or even database is trivial at this point.
func NewServiceForH2(url string) (Service, error) {
	pricesRepo, err := h2Prices.NewRepository(url)
	if err != nil {
		return nil, err
	}

	brandsRepo, err := h2Brands.NewRepository(url)
	if err != nil {
		return nil, err
	}

	return &service{
		pricesRepo: pricesRepo,
		brandsRepo: brandsRepo,
	}, nil
}

func (s *service) Close() {
	if s.pricesRepo != nil {
		s.pricesRepo.Close()
	}
}

func (s *service) ProductPriceForTime(ctx context.Context, brandName string, productID int64, t time.Time) (PriceDetails, error) {
	if s.pricesRepo == nil || s.brandsRepo == nil {
		return PriceDetails{}, fmt.Errorf("no datasource(s) connected")
	}

	brand, err := s.brandsRepo.FindByName(ctx, brandName)
	if err != nil {
		log.Printf("not found brand: %s. Error: %v\n", brandName, err)
		return PriceDetails{}, fmt.Errorf("not found")
	}

	priceLists, err := s.pricesRepo.ListAllForTime(ctx, brand.ID, productID, t)
	if err != nil {
		log.Printf("fetching product price: %v\n", err)
		return PriceDetails{}, fmt.Errorf("error retrieving price")
	}

	if len(priceLists) == 0 {
		return PriceDetails{}, fmt.Errorf("not found")
	}

	p := priceLists[0]
	pd := PriceDetails{
		Price:     p.Price,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		Currency:  p.Currency,
		BrandName: brandName,
	}

	return pd, nil
}
