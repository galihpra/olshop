package handler

type AddressResponse struct {
	Id      uint   `json:"address_id,omitempty"`
	Street  string `json:"street_address,omitempty"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip_code,omitempty"`
}
