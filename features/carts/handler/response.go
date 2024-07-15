package handler

type CartResponse struct {
	Id       uint    `json:"cart_id,omitempty"`
	Quantity int16   `json:"qty,omitempty"`
	Subtotal float64 `json:"subtotal,omitempty"`

	Product ProductResponse `json:"product,omitempty"`
	Varian  VarianResponse  `json:"varian,omitempty"`
}

type ProductResponse struct {
	Id        uint    `json:"product_id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Thumbnail string  `json:"thumbnail,omitempty"`
	Price     float64 `json:"price,omitempty"`
}

type VarianResponse struct {
	Id    uint   `json:"varian_id,omitempty"`
	Color string `json:"color,omitempty"`
}
