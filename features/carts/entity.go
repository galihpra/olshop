package carts

import (
	"context"
	"olshop/helpers/filters"

	"github.com/labstack/echo/v4"
)

type Cart struct {
	ID uint

	ProductID uint
	Product   Product

	VarianID uint
	Varian   Varian

	UserID   uint
	Quantity int16
	Subtotal float64
}

type Product struct {
	ID        uint
	Name      string
	Thumbnail string
	Price     float64
}

type Varian struct {
	ID    uint
	Color string
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, newCart Cart, userId uint) error
	GetAll(ctx context.Context, flt filters.Filter, userId uint) ([]Cart, int, error)
	Update(ctx context.Context, cartId uint, userId uint, updateCart Cart) error
	Delete(ctx context.Context, cartId uint, userId uint) error
}

type Repository interface {
	Create(ctx context.Context, newCart Cart, userId uint) error
	GetAll(ctx context.Context, flt filters.Filter, userId uint) ([]Cart, int, error)
	Update(ctx context.Context, cartId uint, userId uint, updateCart Cart) error
	Delete(ctx context.Context, cartId uint, userId uint) error
}
