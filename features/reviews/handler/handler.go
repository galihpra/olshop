package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/reviews"
	tokens "olshop/helpers/token"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type reviewHandler struct {
	service   reviews.Service
	jwtConfig config.JWT
}

func NewReviewHandler(service reviews.Service, jwtConfig config.JWT) reviews.Handler {
	return &reviewHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (handler *reviewHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(ReviewRequest)
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

		var parseInput = new(reviews.Review)
		parseInput.Review = request.Review
		parseInput.Rating = request.Rating
		parseInput.ProductId = request.ProductId

		if err := handler.service.Create(userId, *parseInput); err != nil {
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

		response["message"] = "create review success"
		return c.JSON(http.StatusCreated, response)
	}
}
