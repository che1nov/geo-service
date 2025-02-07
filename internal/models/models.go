package models

// Запрос на поиск адреса
type SearchRequest struct {
	Query string `json:"query"`
}

// Запрос на геокодирование
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// Модель адреса
type Address struct {
	City string `json:"city"`
}

// Ответ для поиска адреса
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

// Ответ для геокодирования
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

// Запрос на регистрацию пользователя
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Запрос на вход (аутентификацию)
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Ответ при успешном логине
type LoginResponse struct {
	Token string `json:"token"`
}

// Структура для парсинга ответа от DaData API
type DaDataResponse struct {
	Suggestions []struct {
		Data struct {
			City string `json:"city"`
		} `json:"data"`
	} `json:"suggestions"`
}
