package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_service_http_requests_total",
			Help: "Количество HTTP-запросов",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_service_http_request_duration_seconds",
			Help:    "Время выполнения HTTP-запросов",
			Buckets: prometheus.LinearBuckets(0.001, 0.005, 10),
		},
		[]string{"method", "endpoint"},
	)

	RedisRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geo_service_redis_request_duration_seconds",
			Help: "Время выполнения операций с Redis",
		},
		[]string{"operation"},
	)

	DaDataAPIDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geo_service_dadata_api_duration_seconds",
			Help: "Время выполнения запросов к DaData API",
		},
		[]string{"endpoint"},
	)
)
