package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/users"
	tokens "olshop/helpers/token"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	service   users.Service
	jwtConfig config.JWT
}

func NewUserHandler(service users.Service, jwtConfig config.JWT) users.Handler {
	return &userHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (hdl *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(RegisterRequest)
		var response = make(map[string]any)

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "incorect input data"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(users.User)
		parseInput.Name = request.Name
		parseInput.Email = request.Email
		parseInput.Password = request.Password
		parseInput.Username = request.Username

		if err := hdl.service.Register(*parseInput); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "Duplicate") {
				response["message"] = "email is already in use"
				return c.JSON(http.StatusConflict, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "register success"
		return c.JSON(http.StatusCreated, response)
	}
}

func (hdl *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(LoginRequest)
		var response = make(map[string]any)

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "please fill input correctly"
			return c.JSON(http.StatusBadRequest, response)
		}

		result, err := hdl.service.Login(request.Email, request.Password)
		if err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		strToken, err := tokens.GenerateJWT(hdl.jwtConfig.Secret, result.Id)
		if err != nil {
			return err
		}

		var data = new(LoginResponse)
		data.Name = result.Name
		data.Username = result.Username
		data.Email = result.Email
		data.Token = strToken

		response["message"] = "login success"
		response["data"] = data
		return c.JSON(http.StatusOK, response)
	}
}
