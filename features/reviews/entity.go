package reviews

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Review struct {
	ID     uint
	Review string
	Rating float32

	ProductId uint

	User User

	CreatedAt time.Time
}

type User struct {
	Id       uint
	Username string
	Image    string
}

type Product struct {
	Id uint
}

type Handler interface {
	Create() echo.HandlerFunc
}

type Repository interface {
	Create(userId uint, newReview Review) error
}

type Service interface {
	Create(userId uint, newReview Review) error
}
