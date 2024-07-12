package handler

import "time"

type ProductResponse struct {
	Id          uint      `json:"product_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	Discount    int       `json:"discount,omitempty"`
	Rating      float32   `json:"rating,omitempty"`
	Stock       int       `json:"stock,omitempty"`
	Measurement string    `json:"measurement,omitempty"`
	DiscountEnd time.Time `json:"discount_end,omitempty"`

	Thumbnail string   `json:"thumbnail,omitempty"`
	Images    []string `json:"picture,omitempty"`

	Varians []Varianresponse
	Reviews []ReviewResponse
}

type Varianresponse struct {
	Id       uint   `json:"varian_id,omitempty"`
	Color    string `json:"color,omitempty"`
	Stock    int    `json:"stock,omitempty"`
	ImageURL string `json:"image,omitempty"`
}

type ReviewResponse struct {
	ID        uint      `json:"id"`
	Review    string    `json:"text"`
	Rating    float32   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`

	User UserResponse `json:"user"`
}

type UserResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	ImageURL string `json:"image"`
}
