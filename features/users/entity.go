package users

import (
	"context"
	"io"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	Id        uint
	Name      string
	Email     string
	Password  string
	Image     string
	ImageRaw  io.Reader
	Username  string
	CreatedAt time.Time
}

type Handler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Update() echo.HandlerFunc
}

type Service interface {
	Register(newUser User) error
	Login(email string, password string) (*User, error)
	Update(ctx context.Context, id uint, updateUser User) error
}

type Repository interface {
	Register(newUser User) error
	Login(email string) (*User, error)
	Update(ctx context.Context, id uint, updateUser User) error
}
