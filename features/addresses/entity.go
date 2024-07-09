package addresses

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Address struct {
	ID      uint
	Street  string
	Country string
	City    string
	State   string
	Zip     string

	UserID uint
	User   User
}

type User struct {
	ID uint
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, data Address) error
	GetAll(ctx context.Context) ([]Address, error)
	Delete(ctx context.Context, id uint) error
}

type Repository interface {
	Create(ctx context.Context, data Address) error
	GetAll(ctx context.Context) ([]Address, error)
	Delete(ctx context.Context, id uint) error
}
