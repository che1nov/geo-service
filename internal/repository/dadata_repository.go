package repository

import (
	"bytes"
	"encoding/json"
	"net/http"

	"example/internal/models"
)

// DaDataRepository определяет интерфейс для работы с DaData API.
type DaDataRepository interface {
	CallDaDataAPI(endpoint string, payload interface{}) (*models.DaDataResponse, error)
}

type DataRepository struct {
	apiURL string
	apiKey string
}

// NewDaDataRepository возвращает реализацию DaDataRepository.
func NewDataRepository() DaDataRepository {
	return &DataRepository{
		apiURL: "https://suggestions.dadata.ru/suggestions/api/4_1/rs",
		apiKey: "2978b2a44fd0275be5a33f228be6d4bdff6aa1f9",
	}
}

func (r *DataRepository) CallDaDataAPI(endpoint string, payload interface{}) (*models.DaDataResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := r.apiURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+r.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.DaDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
