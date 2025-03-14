package middleware

import (
	"geo-service/internal/metrics"
	"net/http"
	"strconv"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для ResponseWriter чтобы отслеживать статус код
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполняем запрос
		next.ServeHTTP(ww, r)

		// Измеряем время выполнения
		duration := time.Since(start).Seconds()

		// Обновляем метрики
		metrics.HttpRequestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.statusCode)).Inc()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
