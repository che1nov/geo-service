package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"geo-service/internal/entities"
	"geo-service/internal/repository"
)

type AddressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) Search(ctx context.Context, query string) (*entities.ResponseAddress, error) {
	addresses, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	entityAddresses := make([]*entities.Address, 0, len(addresses))
	for _, addr := range addresses {
		entityAddr := MapToEntityAddress(addr)
		if entityAddr != nil {
			entityAddresses = append(entityAddresses, entityAddr)
		}
	}

	return &entities.ResponseAddress{Addresses: entityAddresses}, nil
}

func (s *AddressService) Geocode(ctx context.Context, lat, lng float64) (*entities.ResponseAddress, error) {
	addresses, err := s.repo.Geocode(ctx, lat, lng)
	if err != nil {
		return nil, err
	}

	entityAddresses := make([]*entities.Address, 0, len(addresses))
	for _, addr := range addresses {
		entityAddr := MapToEntityAddress(addr)
		if entityAddr != nil {
			entityAddresses = append(entityAddresses, entityAddr)
		}
	}

	return &entities.ResponseAddress{Addresses: entityAddresses}, nil
}

func MapToEntityAddress(repoAddr *repository.DaDataAddress) *entities.Address {
	if repoAddr == nil || repoAddr.Data == nil {
		logrus.Warn("Попытка преобразования nil адреса")
		return nil
	}

	getString := func(data map[string]interface{}, key string) string {
		if val, ok := data[key].(string); ok {
			return val
		}
		logrus.WithField("key", key).Warn("Пропущен ключ в данных DaData")
		return ""
	}

	return &entities.Address{
		City:    getString(repoAddr.Data, "city"),
		Street:  getString(repoAddr.Data, "street"),
		House:   getString(repoAddr.Data, "house"),
		ZipCode: getString(repoAddr.Data, "postal_code"),
	}
}
