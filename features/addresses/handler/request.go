package handler

type AddressRequest struct {
	Street  string `json:"street_address"`
	City    string `json:"city"`
	Country string `json:"country"`
	State   string `json:"state"`
	Zip     string `json:"zip_code"`
}
