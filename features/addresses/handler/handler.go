package handler

import (
	"context"
	"net/http"
	"olshop/config"
	"olshop/features/addresses"
	tokens "olshop/helpers/token"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type addressHandler struct {
	service   addresses.Service
	jwtConfig config.JWT
}

func NewAddressHandler(service addresses.Service, jwtConfig config.JWT) addresses.Handler {
	return &addressHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (handler *addressHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(AddressRequest)
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

		var parseInput = new(addresses.Address)
		parseInput.Street = request.Street
		parseInput.City = request.City
		parseInput.Country = request.Country
		parseInput.State = request.State
		parseInput.Zip = request.Zip
		parseInput.UserID = userId

		if err := handler.service.Create(c.Request().Context(), *parseInput); err != nil {
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

		response["message"] = "create address success"
		return c.JSON(http.StatusCreated, response)
	}
}

func (handler *addressHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid product id"
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

		if err := handler.service.Delete(c.Request().Context(), uint(id), userId); err != nil {
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

		response["message"] = "delete address success"
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *addressHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		result, err := handler.service.GetAll(context.Background())
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data []AddressResponse
		for _, address := range result {
			data = append(data, AddressResponse{
				Id:      address.ID,
				Street:  address.Street,
				City:    address.City,
				Country: address.Country,
				State:   address.State,
				Zip:     address.Zip,
			})
		}

		response["message"] = "get all address success"
		response["data"] = data
		return c.JSON(http.StatusOK, response)
	}
}
