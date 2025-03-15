package handlers

import (
	"encoding/json"
	"net/http"

	"geo-service/internal/entities"
	"geo-service/internal/service"
	"github.com/sirupsen/logrus"
)

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
