package mux

import (
	"kafka-connect/internal/infra/http/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMux(mh *handlers.MetricsHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(5, "text/html"))
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/metrics", mh.Handler)

	return r
}
