package server

import "github.com/mehix/pricelists/service/pricelist"

type Option func(*server)

func WithDebug(dbg bool) Option {
	return func(s *server) {
		s.debug = dbg
	}
}

func WithPricelist(svc pricelist.Service) Option {
	return func(s *server) {
		s.pricelistSvc = svc
	}
}
