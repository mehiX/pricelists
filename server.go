package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mehix/pricelists/service/pricelist"
)

type server struct {
	debug        bool
	pricelistSvc pricelist.Service
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

	r.Get("/prices/prod/{productID}/brand/{brandName}/date/{date:[0-9-]{10}}/time/{time:[0-9:]{8}}", s.handlePriceRequest)

	return r
}

type PriceResponse struct {
	Price        int64  `json:"price"`
	PriceDisplay string `json:"price_display"`
}

func (s *server) handlePriceRequest(w http.ResponseWriter, r *http.Request) {
	if s.pricelistSvc == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-type", "application/json")

	// read and validate request parameters
	productID, err := strconv.ParseInt(chi.URLParam(r, "productID"), 10, 64)
	if err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}
	brandName := chi.URLParam(r, "brandName")
	d := chi.URLParam(r, "date")
	t := chi.URLParam(r, "time")

	datetime, err := time.Parse("2006-01-02 15:04:05", d+" "+t)
	if err != nil {
		respondError(w, err, http.StatusBadRequest)
		return
	}

	// call business logic
	priceDetails, err := s.pricelistSvc.ProductPriceForTime(r.Context(), brandName, productID, datetime)
	if err != nil {
		respondError(w, err, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(PriceResponse{
		Price:        priceDetails.Price,
		PriceDisplay: fmt.Sprintf("%.2f", float64(priceDetails.Price)/100),
	})
}

func respondError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct{ Error string }{Error: err.Error()})
}
