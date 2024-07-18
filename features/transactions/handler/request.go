package handler

type CartItem struct {
	CartID uint `json:"cart_id"`
}

type TransactionRequest struct {
	AddressId     uint       `json:"address_id"`
	PaymentMethod string     `json:"payment_method"`
	Items         []CartItem `json:"items"`
}
