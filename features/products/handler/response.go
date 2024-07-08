package handler

type ProductResponse struct {
	Id          uint    `json:"product_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Discount    int     `json:"discount,omitempty"`
	Rating      float32 `json:"rating,omitempty"`

	Thumbnail string   `json:"thumbnail,omitempty"`
	Images    []string `json:"images,omitempty"`
}
