package product

import (
	"context"
	"io"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ID        uint
	Name      string
	Rating    float32
	Price     float64
	Thumbnail Image

	Images []Image
}

type Image struct {
	ID       uint
	ImageURL string
	ImageRaw io.Reader
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, data Product) error
	GetAll(ctx context.Context) ([]Product, int, error)
}

type Repository interface {
	Create(ctx context.Context, data Product) error
	GetAll(ctx context.Context) ([]Product, int, error)
}
