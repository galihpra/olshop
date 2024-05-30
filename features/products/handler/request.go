package handler

import "io"

type CreateRequest struct {
	Name   string  `json:"name" form:"name"`
	Price  float64 `json:"price" form:"price"`
	Images []io.Reader

	CategoryId uint `formam:"category_id"`
}
