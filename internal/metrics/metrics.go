package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP метрики
var (
	// Количество запросов по эндпоинтам, методам и статусам
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_service_http_requests_total",
			Help: "Количество HTTP-запросов по эндпоинтам, методам и статусам",
		},
		[]string{"endpoint", "method", "status"},
	)

	// Время выполнения HTTP-запросов
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_service_http_request_duration_seconds",
			Help:    "Время выполнения HTTP-запросов",
			Buckets: prometheus.LinearBuckets(0.001, 0.005, 10),
		},
		[]string{"endpoint", "method"},
	)
)

// Redis метрики
var (
	RedisRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geo_service_redis_request_duration_seconds",
			Help: "Время выполнения операций с Redis",
		},
		[]string{"operation"},
	)
)

// База данных метрики
var (
	DbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geo_service_db_query_duration_seconds",
			Help: "Время выполнения SQL-запросов",
		},
		[]string{"query_type"},
	)
)

// DaData API метрики
var (
	DaDataAPIDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geo_service_dadata_api_duration_seconds",
			Help: "Время выполнения запросов к DaData API",
		},
		[]string{"api_endpoint"},
	)

	DaDataAPIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_service_dadata_api_requests_total",
			Help: "Количество запросов к DaData API",
		},
		[]string{"endpoint", "status"},
	)
)

// Аутентификационные метрики
var (
	AuthRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_service_auth_requests_total",
			Help: "Количество запросов к аутентификационным эндпоинтам",
		},
		[]string{"endpoint", "status"},
	)
)
