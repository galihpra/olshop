package products

import (
	"context"
	"io"
	"olshop/helpers/filters"
	"time"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ID          uint
	Name        string
	Rating      float32
	Price       float64
	Thumbnail   string
	Discount    int
	Description string
	Stock       int
	DiscountEnd time.Time
	Measurement string

	Images   []Image
	Category Category
	Varians  []Varian
}

type Image struct {
	ID       uint
	ImageURL string
	ImageRaw io.Reader
}

type Category struct {
	ID       uint
	Category string
}

type Varian struct {
	ID       uint
	Color    string
	Stock    int
	ImageURL string
	ImageRaw io.Reader
}

type Handler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetProductDetail() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Create(ctx context.Context, data Product) error
	GetAll(ctx context.Context, flt filters.Filter) ([]Product, int, error)
	Update(ctx context.Context, updateProduct Product) error
	GetProductDetail(ctx context.Context, id uint) (*Product, error)
	Delete(ctx context.Context, id uint) error
}

type Repository interface {
	Create(ctx context.Context, data Product) error
	GetAll(ctx context.Context, flt filters.Filter) ([]Product, int, error)
	Update(ctx context.Context, updateProduct Product) error
	GetProductDetail(ctx context.Context, id uint) (*Product, error)
	Delete(ctx context.Context, id uint) error
}
