package users

import (
	"time"

	"github.com/labstack/echo"
)

type User struct {
	Id        uint
	Name      string
	Email     string
	Password  string
	Image     string
	Username  string
	CreatedAt time.Time
}

type Handler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
}

type Service interface {
	Register(newUser User) error
	Login(email string, password string) (*User, error)
}

type Repository interface {
	Register(newUser User) error
	Login(email string) (*User, error)
}
