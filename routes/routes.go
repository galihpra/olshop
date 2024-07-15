package routes

import (
	"olshop/features/addresses"
	"olshop/features/carts"
	"olshop/features/products"
	"olshop/features/reviews"
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
	ReviewHandler  reviews.Handler
	CartHandler    carts.Handler
}

func (router Routes) InitRouter() {
	router.UserRouter()
	router.ProductRouter()
	router.AddressRouter()
	router.ReviewRouter()
	router.CartRouter()
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
	router.Server.PUT("/products/:id", router.ProductHandler.Update(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/products/:id/reviews", router.ProductHandler.GetAllReview())
}

func (router *Routes) AddressRouter() {
	router.Server.POST("/addresses", router.AddressHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.GET("/addresses", router.AddressHandler.GetAll(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.DELETE("/addresses/:id", router.AddressHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
}

func (router *Routes) ReviewRouter() {
	router.Server.POST("/reviews", router.ReviewHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
}
func (router *Routes) CartRouter() {
	router.Server.POST("/carts", router.CartHandler.Create(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.DELETE("/carts/:id", router.CartHandler.Delete(), echojwt.JWT([]byte(router.JWTKey)))
	router.Server.PUT("/carts/:id", router.CartHandler.Update(), echojwt.JWT([]byte(router.JWTKey)))
}
