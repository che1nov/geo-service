package router

import (
	"geo-service/internal/handlers"
	"geo-service/internal/metrics"
	"geo-service/internal/middleware"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
)

func NewChiRouter(addressHandler *handlers.AddressHandler, authHandler *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	r.Post("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		metrics.AuthRequestsTotal.WithLabelValues("/api/auth/register", "200").Inc()
		authHandler.Register(w, r)
	})

	r.Post("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		metrics.AuthRequestsTotal.WithLabelValues("/api/auth/login", "200").Inc()
		authHandler.Login(w, r)
	})

	protected := r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
		r.Post("/address/search", func(w http.ResponseWriter, r *http.Request) {
			metrics.HttpRequestsTotal.WithLabelValues("/api/address/search", "POST", "200").Inc()
			addressHandler.Search(w, r)
		})
		r.Post("/address/geocode", func(w http.ResponseWriter, r *http.Request) {
			metrics.HttpRequestsTotal.WithLabelValues("/api/address/geocode", "POST", "200").Inc()
			addressHandler.Geocode(w, r)
		})
	})

	r.Use(middleware.MetricsMiddleware)

	return r
}
