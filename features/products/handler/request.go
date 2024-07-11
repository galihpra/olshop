package handler

import (
	"io"
	"time"
)

type CreateRequest struct {
	Name        string  `json:"name" form:"name"`
	Price       float64 `json:"price" form:"price"`
	Images      []io.Reader
	CategoryId  uint      `form:"category_id"`
	Discount    int       `json:"discount" form:"discount"`
	Stock       int       `json:"stock" form:"stock"`
	Description string    `json:"description" form:"description"`
	Measurement string    `json:"measurement" form:"measurement"`
	DiscountEnd time.Time `json:"discount_end" form:"discount_end"`

	Varians []VarianRequest
}

type VarianRequest struct {
	Color    string    `json:"color" form:"color"`
	Stock    int       `json:"stock" form:"stock"`
	ImageRaw io.Reader `json:"varian_image" form:"varian_image"`
}
