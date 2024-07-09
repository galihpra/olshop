package routes

import (
	"olshop/features/addresses"
	"olshop/features/products"
	"olshop/features/users"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	JWTKey         string
	Server         *echo.Echo
	UserHandler    users.Handler
	ProductHandler products.Handler
	AddressHandler addresses.Handler
}

func (router Routes) InitRouter() {
	router.UserRouter()
	router.ProductRouter()
	router.AddressRouter()
}

func (router *Routes) UserRouter() {
	router.Server.POST("/register", router.UserHandler.Register())
	router.Server.POST("/login", router.UserHandler.Login())
	router.Server.PATCH("/users", router.UserHandler.Update(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.DELETE("/users", router.UserHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/users", router.UserHandler.GetById(), echojwt.JWT([]byte(router.JWTKey)))
}

func (router *Routes) ProductRouter() {
	router.Server.POST("/products", router.ProductHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/products", router.ProductHandler.GetAll())
	router.Server.GET("/products/:id", router.ProductHandler.GetProductDetail())
	router.Server.DELETE("/products/:id", router.ProductHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
}

func (router *Routes) AddressRouter() {
	router.Server.POST("/addresses", router.AddressHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
}
