package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/carts"
	tokens "olshop/helpers/token"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type cartHandler struct {
	service   carts.Service
	jwtConfig config.JWT
}

func NewCartHandler(service carts.Service, jwtConfig config.JWT) carts.Handler {
	return &cartHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (handler *cartHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(CartRequest)
		var response = make(map[string]any)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(handler.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "incorect input data"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(carts.Cart)
		parseInput.ProductID = request.ProductId
		parseInput.VarianID = request.VarianId
		parseInput.Quantity = request.Quantity
		parseInput.UserID = userId

		if err := handler.service.Create(c.Request().Context(), *parseInput, uint(userId)); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "unauthorized") {
				response["message"] = "unauthorized"
				return c.JSON(http.StatusBadRequest, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "product added to cart"
		return c.JSON(http.StatusCreated, response)
	}
}

func (handler *cartHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		cartId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid cart id"
		}

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(handler.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := handler.service.Delete(c.Request().Context(), uint(cartId), userId); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "not found") {
				response["message"] = "not found"
				return c.JSON(http.StatusNotFound, response)
			}

			if strings.Contains(err.Error(), "invalid id") {
				response["message"] = "not found"
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "delete cart success"
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *cartHandler) GetAll() echo.HandlerFunc {
	panic("unimplemented")
}

func (handler *cartHandler) Update() echo.HandlerFunc {
	panic("unimplemented")
}
