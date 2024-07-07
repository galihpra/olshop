package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/users"
	tokens "olshop/helpers/token"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
		data.Image = result.Image
		data.Token = strToken

		response["message"] = "login success"
		response["data"] = data
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *userHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var request = new(RegisterRequest)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(hdl.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "please fill input correctly"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(users.User)
		parseInput.Name = request.Name
		parseInput.Email = request.Email
		parseInput.Password = request.Password
		parseInput.Username = request.Username

		file, _ := c.FormFile("image")
		if file != nil {
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			request.Image = src
		}
		parseInput.ImageRaw = request.Image

		if err := hdl.service.Update(userId, *parseInput); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "validate: ") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}

			if strings.Contains(err.Error(), "not found: ") {
				response["message"] = "user not found"
				return c.JSON(http.StatusNotFound, response)
			}

			if strings.Contains(err.Error(), "Duplicate") {
				response["message"] = "this email has been used, please use another email"
				return c.JSON(http.StatusConflict, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "update user success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *userHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(hdl.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := hdl.service.Delete(userId); err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "not found: ") {
				response["message"] = "user not found"
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		response["message"] = "delete user success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *userHandler) GetById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		userId, err := tokens.ExtractToken(hdl.jwtConfig.Secret, token.(*jwt.Token))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "unauthorized"
			return c.JSON(http.StatusUnauthorized, response)
		}

		result, err := hdl.service.GetById(userId)
		if err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "not found: ") {
				response["message"] = "user not found"
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data = new(UserResponse)
		data.Name = result.Name
		data.Username = result.Username
		data.Email = result.Email
		data.Image = result.Image

		response["message"] = "get user success"
		response["data"] = data
		return c.JSON(http.StatusOK, response)

	}
}
