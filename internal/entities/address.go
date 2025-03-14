package entities

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	City    string `json:"city"`
	Street  string `json:"street"`
	House   string `json:"house"`
	ZipCode string `json:"zip_code"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}
