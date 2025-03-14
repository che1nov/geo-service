package handlers

import (
	"encoding/json"
	"net/http"

	"geo-service/internal/entities"
	"geo-service/internal/service"
	"github.com/sirupsen/logrus"
)

// swagger:meta
// GeoService API предоставляет функциональность для поиска адресов и геокодирования.
// BasePath: /api/address
type AddressHandler struct {
	service *service.AddressService
	log     *logrus.Logger
}

func NewAddressHandler(service *service.AddressService, log *logrus.Logger) *AddressHandler {
	return &AddressHandler{
		service: service,
		log:     log,
	}
}

// swagger:route POST /search Address searchAddress
// Поиск адреса по строковому запросу.
// responses:
//
//	200: addressResponse
//	400: errorResponse
//	500: errorResponse
//
// swagger:parameters searchAddress
func (h *AddressHandler) Search(w http.ResponseWriter, r *http.Request) {
	var request entities.RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.Search(r.Context(), request.Query)
	if err != nil {
		http.Error(w, "Failed to search addresses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /geocode Address geocodeAddress
// Геокодирование координат в адрес.
// responses:
//
//	200: addressResponse
//	400: errorResponse
//	500: errorResponse
//
// swagger:parameters geocodeAddress
func (h *AddressHandler) Geocode(w http.ResponseWriter, r *http.Request) {
	var request entities.RequestAddressGeocode
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.Geocode(r.Context(), request.Lat, request.Lng)
	if err != nil {
		http.Error(w, "Failed to geocode addresses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
