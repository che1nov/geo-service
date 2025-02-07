package controller

import (
	"encoding/json"
	"net/http"

	"example/internal/models"
	"example/internal/service"
)

// Responder интерфейс для формирования ответа.
type Responder interface {
	RespondJSON(w http.ResponseWriter, status int, data interface{})
	RespondError(w http.ResponseWriter, status int, message string)
}

type AddressController struct {
	addressService service.AddressService
	responder      Responder
}

func NewAddressController(addressService service.AddressService, responder Responder) *AddressController {
	return &AddressController{
		addressService: addressService,
		responder:      responder,
	}
}

// Search обрабатывает поиск адреса через DaData API.
func (c *AddressController) Search(w http.ResponseWriter, r *http.Request) {
	var req models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := c.addressService.Search(req.Query)
	if err != nil {
		c.responder.RespondError(w, http.StatusInternalServerError, "Failed to search address")
		return
	}

	c.responder.RespondJSON(w, http.StatusOK, response)
}

// Geocode обрабатывает геокодирование адреса через DaData API.
func (c *AddressController) Geocode(w http.ResponseWriter, r *http.Request) {
	var req models.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := c.addressService.Geocode(req.Lat, req.Lng)
	if err != nil {
		c.responder.RespondError(w, http.StatusInternalServerError, "Failed to geocode address")
		return
	}

	c.responder.RespondJSON(w, http.StatusOK, response)
}
