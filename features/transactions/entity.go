package transactions

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type Transaction struct {
	Invoice         int
	Total           float64
	PaymentMethod   string
	Status          string
	TransactionDate time.Time

	UserID uint
	User   User

	AddressID uint
	Address   Address

	TransactionDetails []TransactionDetail

	Payment Payment
}

type User struct {
	ID    uint
	Name  string
	Email string
}

type TransactionDetail struct {
	ID       uint
	Quantity int16
	Subtotal float64

	Products Product
}

type Address struct {
	ID uint
}

type Product struct {
	ID    uint
	Name  string
	Price float64
}

type Payment struct {
	Id            uint
	Method        string
	Bank          string
	VirtualNumber string
	BillKey       string
	BillCode      string
	Status        string

	TransactionCode  int
	TransactionTotal float64

	CreatedAt time.Time
	ExpiredAt time.Time
	PaidAt    time.Time
}

type Handler interface {
	Create() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, userId uint, cartIds []uint, newTransaction Transaction) (*Transaction, error)
}

type Repository interface {
	Create(ctx context.Context, userId uint, cartIds []uint, newTransaction Transaction) (*Transaction, error)
}
