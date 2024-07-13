package carts

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Cart struct {
	ID        uint
	ProductID uint
	VarianID  uint
	UserID    uint
	Quantity  int16
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, newCart Cart, userId uint) error
	GetAll(ctx context.Context, userId uint) ([]Cart, error)
	Update(ctx context.Context, cartId uint, userId uint, updateCart Cart) error
	Delete(ctx context.Context, cartId uint, userId uint) error
}

type Repository interface {
	Create(ctx context.Context, newCart Cart, userId uint) error
	GetAll(ctx context.Context, userId uint) ([]Cart, error)
	Update(ctx context.Context, cartId uint, userId uint, updateCart Cart) error
	Delete(ctx context.Context, cartId uint, userId uint) error
}
