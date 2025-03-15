package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"

	"geo-service/internal/handlers"
	"geo-service/internal/middleware"
)

func NewChiRouter(addressHandler *handlers.AddressHandler, authHandler *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
		r.Use(middleware.MetricsMiddleware)

		r.Post("/address/search", addressHandler.Search)
		r.Post("/address/geocode", addressHandler.Geocode)
	})

	return r
}
