package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"geo-service/internal/metrics"
)

type AddressRepository interface {
	Search(ctx context.Context, query string) ([]*DaDataAddress, error)
	Geocode(ctx context.Context, lat, lng float64) ([]*DaDataAddress, error)
}

type DaDataRepository struct {
	apiKey string
	apiURL string
	log    *logrus.Logger
	cache  *redis.Client
}

type DaDataAddress struct {
	Value           string                 `json:"value"`
	UnrestrictedVal string                 `json:"unrestricted_value"`
	Data            map[string]interface{} `json:"data"`
}

func NewDaDataRepository(apiKey, apiURL string, log *logrus.Logger, redisClient *redis.Client) AddressRepository {
	return &DaDataRepository{
		apiKey: apiKey,
		apiURL: apiURL,
		log:    log,
		cache:  redisClient,
	}
}

func (r *DaDataRepository) Search(ctx context.Context, query string) ([]*DaDataAddress, error) {
	cacheKey := fmt.Sprintf("search:%s", query)

	start := time.Now()
	cachedData, err := r.cache.Get(ctx, cacheKey).Result()
	metrics.RedisRequestDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
	if err == nil {
		r.log.WithField("cache_key", cacheKey).Info("Данные получены из кэша")
		var addresses []*DaDataAddress
		if err := json.Unmarshal([]byte(cachedData), &addresses); err != nil {
			r.log.WithError(err).Warn("Ошибка десериализации данных из кэша")
			// Продолжаем запрос к API
		}
	}

	payload := map[string]string{"query": query}
	data, err := r.callAPI(ctx, "/suggest/address", payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка вызова API: %w", err)
	}

	var resp struct {
		Suggestions []*DaDataAddress `json:"suggestions"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %w", err)
	}

	jsonData, err := json.Marshal(resp.Suggestions)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации данных для кэша: %w", err)
	}

	startSet := time.Now()
	if err := r.cache.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
		metrics.RedisRequestDuration.WithLabelValues("SET").Observe(time.Since(startSet).Seconds())
		r.log.WithError(err).Warnf("Ошибка записи в кэш: %s", cacheKey)
	}

	return resp.Suggestions, nil
}

func (r *DaDataRepository) Geocode(ctx context.Context, lat, lng float64) ([]*DaDataAddress, error) {
	cacheKey := fmt.Sprintf("geocode:%.6f:%.6f", lat, lng)

	start := time.Now()
	cachedData, err := r.cache.Get(ctx, cacheKey).Result()
	metrics.RedisRequestDuration.WithLabelValues("GET").Observe(time.Since(start).Seconds())
	if err == nil {
		r.log.WithField("cache_key", cacheKey).Info("Данные получены из кэша")
		var addresses []*DaDataAddress
		if err := json.Unmarshal([]byte(cachedData), &addresses); err != nil {
			r.log.WithError(err).Warnf("Ошибка десериализации данных из кэша: %s", cacheKey)
		}
	}

	payload := map[string]string{
		"lat": strconv.FormatFloat(lat, 'f', 6, 64),
		"lon": strconv.FormatFloat(lng, 'f', 6, 64),
	}

	data, err := r.callAPI(ctx, "/geolocate/address", payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка вызова API: %w", err)
	}

	var resp struct {
		Suggestions []*DaDataAddress `json:"suggestions"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %w", err)
	}

	jsonData, err := json.Marshal(resp.Suggestions)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации данных для кэша: %w", err)
	}

	// Запись в кэш
	startSet := time.Now()
	if err := r.cache.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
		metrics.RedisRequestDuration.WithLabelValues("SET").Observe(time.Since(startSet).Seconds())
		r.log.WithError(err).Warnf("Ошибка записи в кэш: %s", cacheKey)
	}

	return resp.Suggestions, nil
}

func (r *DaDataRepository) callAPI(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	url := r.apiURL + endpoint
	body, err := json.Marshal(payload)
	if err != nil {
		r.log.WithError(err).Error("Ошибка сериализации payload")
		metrics.DaDataAPIDuration.WithLabelValues(endpoint).Observe(0)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		r.log.WithError(err).Error("Ошибка создания HTTP-запроса")
		metrics.DaDataAPIDuration.WithLabelValues(endpoint).Observe(0)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+r.apiKey)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Do(req)
	metrics.DaDataAPIDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds())
	if err != nil {
		r.log.WithError(err).Error("Ошибка выполнения HTTP-запроса к DaData")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		r.log.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(bodyBytes),
		}).Error("Ошибка DaData API")
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.WithError(err).Error("Ошибка чтения ответа от DaData")
		return nil, err
	}

	return data, nil
}
