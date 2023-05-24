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

	return r
}
