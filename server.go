package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type server struct {
	debug bool
}

func New(opts ...Option) *server {
	s := &server{}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *server) Handlers() http.Handler {

	r := chi.NewMux()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/prices/prod/{productID}/brand/{brandName}/date/{date:[0-9-]{10}}/time/{time:[0-9:]{8}}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	return r
}

type PriceResponse struct {
	Price        int64  `json:"price"`
	PriceDisplay string `json:"price_display"`
}
