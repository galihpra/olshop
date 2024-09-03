package handler

import "time"

type TransactionResponse struct {
	Invoice int     `json:"invoice,omitempty"`
	Total   float64 `json:"total,omitempty"`
	Status  string  `json:"status,omitempty"`

	PaymentBank          string     `json:"payment_method,omitempty"`
	PaymentVirtualNumber string     `json:"virtual_number,omitempty"`
	PaymentBillKey       string     `json:"key_bill,omitempty"`
	PaymentBillCode      string     `json:"code_bill,omitempty"`
	PaymentExpiredAt     *time.Time `json:"payment_expired,omitempty"`
}
