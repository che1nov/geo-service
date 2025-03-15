package middleware

import (
	"fmt"
	"geo-service/internal/metrics"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		method := r.Method
		endpoint := r.URL.Path
		status := recorder.status

		metrics.HttpRequestsTotal.WithLabelValues(method, endpoint, fmt.Sprint(status)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	})
}
