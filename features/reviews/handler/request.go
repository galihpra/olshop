package handler

type ReviewRequest struct {
	Review    string  `json:"text"`
	Rating    float32 `json:"rating"`
	ProductId uint    `json:"product_id"`
}
