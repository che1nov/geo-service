package service

import (
	"example/internal/models"
	"example/internal/repository"
)

// AddressService описывает бизнес-логику для работы с адресами.
type AddressService interface {
	Search(query string) (*models.SearchResponse, error)
	Geocode(lat, lng string) (*models.GeocodeResponse, error)
}

type addressService struct {
	dadataRepo repository.DaDataRepository
}

func NewAddressService(dadataRepo repository.DaDataRepository) AddressService {
	return &addressService{
		dadataRepo: dadataRepo,
	}
}

// Search godoc
// @Summary      Поиск адреса
// @Description  Выполняет поиск адреса по запросу через DaData API и возвращает найденные адреса.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        searchRequest  body      models.SearchRequest  true  "Запрос для поиска адреса"
// @Success      200            {object}  models.SearchResponse  "Найденные адреса"
// @Failure      400            {object}  map[string]string  "Неверный запрос"
// @Failure      500            {object}  map[string]string  "Ошибка при вызове DaData API"
// @Security     ApiKeyAuth
// @Router       /api/address/search [post]
func (s *addressService) Search(query string) (*models.SearchResponse, error) {
	urlEndpoint := "/suggest/address"
	payload := map[string]string{"query": query}
	response, err := s.dadataRepo.CallDaDataAPI(urlEndpoint, payload)
	if err != nil {
		return nil, err
	}

	result := &models.SearchResponse{}
	for _, suggestion := range response.Suggestions {
		result.Addresses = append(result.Addresses, &models.Address{City: suggestion.Data.City})
	}
	return result, nil
}

// Geocode godoc
// @Summary      Геокодирование адреса
// @Description  Выполняет геокодирование адреса по координатам через DaData API и возвращает найденные адреса.
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        geocodeRequest  body      models.GeocodeRequest  true  "Запрос для геокодирования (lat, lng)"
// @Success      200             {object}  models.GeocodeResponse  "Найденные адреса"
// @Failure      400             {object}  map[string]string  "Неверный запрос"
// @Failure      500             {object}  map[string]string  "Ошибка при вызове DaData API"
// @Security     ApiKeyAuth
// @Router       /api/address/geocode [post]
func (s *addressService) Geocode(lat, lng string) (*models.GeocodeResponse, error) {
	urlEndpoint := "/geolocate/address"
	payload := map[string]string{"lat": lat, "lon": lng}
	response, err := s.dadataRepo.CallDaDataAPI(urlEndpoint, payload)
	if err != nil {
		return nil, err
	}

	result := &models.GeocodeResponse{}
	for _, suggestion := range response.Suggestions {
		result.Addresses = append(result.Addresses, &models.Address{City: suggestion.Data.City})
	}
	return result, nil
}
