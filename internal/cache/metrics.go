package cache

import (
	"geo-service/internal/metrics"
	"time"
)

func measureCacheOperation(operation string, f func() error) error {
	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()
	metrics.RedisRequestDuration.WithLabelValues(operation).Observe(duration)
	return err
}

func measureCacheOperationWithResult[T any](operation string, f func() (T, error)) (T, error) {
	start := time.Now()
	result, err := f()
	duration := time.Since(start).Seconds()
	metrics.RedisRequestDuration.WithLabelValues(operation).Observe(duration)
	return result, err
}
