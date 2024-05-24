package routes

import (
	"olshop/features/users"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	JWTKey      string
	Server      *echo.Echo
	UserHandler users.Handler
}

func (router Routes) InitRouter() {
	router.UserRouter()
}

func (router *Routes) UserRouter() {
	router.Server.POST("/register", router.UserHandler.Register())
	router.Server.POST("/login", router.UserHandler.Login())
}
