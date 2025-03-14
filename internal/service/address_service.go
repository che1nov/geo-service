package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"geo-service/internal/entities"
	"geo-service/internal/metrics"
	"geo-service/internal/repository"
)

type AddressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) Search(ctx context.Context, query string) (*entities.ResponseAddress, error) {
	start := time.Now()
	defer func() {
		metrics.AddressServiceSearchDuration.Observe(time.Since(start).Seconds())
	}()

	addresses, err := s.repo.Search(ctx, query)
	if err != nil {
		metrics.AddressServiceSearchRequestsTotal.WithLabelValues("error").Inc()
		logrus.WithError(err).Error("Ошибка при поиске адреса")
		return nil, err
	}

	metrics.AddressServiceSearchRequestsTotal.WithLabelValues("success").Inc()

	entityAddresses := make([]*entities.Address, 0, len(addresses))
	for _, addr := range addresses {
		entityAddr := MapToEntityAddress(addr)
		if entityAddr == nil {
			logrus.Warn("Пропущен адрес с nil данными")
			continue
		}
		entityAddresses = append(entityAddresses, entityAddr)
	}

	return &entities.ResponseAddress{
		Addresses: entityAddresses,
	}, nil
}

func (s *AddressService) Geocode(ctx context.Context, lat, lng float64) (*entities.ResponseAddress, error) {
	start := time.Now()
	defer func() {
		metrics.AddressServiceGeocodeDuration.Observe(time.Since(start).Seconds())
	}()

	addresses, err := s.repo.Geocode(ctx, lat, lng)
	if err != nil {
		metrics.AddressServiceGeocodeRequestsTotal.WithLabelValues("error").Inc()
		logrus.WithError(err).Error("Ошибка при геокодировании")
		return nil, err
	}

	metrics.AddressServiceGeocodeRequestsTotal.WithLabelValues("success").Inc()

	entityAddresses := make([]*entities.Address, 0, len(addresses))
	for _, addr := range addresses {
		entityAddr := MapToEntityAddress(addr)
		if entityAddr == nil {
			logrus.Warn("Пропущен адрес с nil данными")
			continue
		}
		entityAddresses = append(entityAddresses, entityAddr)
	}

	return &entities.ResponseAddress{
		Addresses: entityAddresses,
	}, nil
}

func MapToEntityAddress(repoAddr *repository.DaDataAddress) *entities.Address {
	if repoAddr == nil || repoAddr.Data == nil {
		return nil
	}
	getString := func(data map[string]interface{}, key string) string {
		if val, ok := data[key].(string); ok {
			return val
		}
		return ""
	}

	return &entities.Address{
		City:    getString(repoAddr.Data, "city"),
		Street:  getString(repoAddr.Data, "street"),
		House:   getString(repoAddr.Data, "house"),
		ZipCode: getString(repoAddr.Data, "postal_code"),
	}
}
