package handler

import (
	"context"
	"fmt"
	"net/http"
	"olshop/config"
	"olshop/features/carts"
	"olshop/helpers/filters"
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
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var baseUrl = c.Scheme() + "://" + c.Request().Host

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

		var pagination = new(filters.Pagination)
		c.Bind(pagination)
		if pagination.Start != 0 && pagination.Limit == 0 {
			pagination.Limit = 5
		}

		var search = new(filters.Search)
		c.Bind(search)

		var sort = new(filters.Sort)
		c.Bind(sort)

		result, totalData, err := handler.service.GetAll(context.Background(), filters.Filter{Search: *search, Pagination: *pagination, Sort: *sort}, userId)
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data []CartResponse
		for _, cart := range result {
			data = append(data, CartResponse{
				Id:       cart.ID,
				Quantity: cart.Quantity,
				Subtotal: cart.Subtotal,
				Product: ProductResponse{
					Id:        cart.Product.ID,
					Name:      cart.Product.Name,
					Price:     cart.Product.Price,
					Thumbnail: cart.Product.Thumbnail,
				},
				Varian: VarianResponse{
					Id:    cart.Varian.ID,
					Color: cart.Varian.Color,
				},
			})
		}
		response["data"] = data

		if pagination.Limit != 0 {
			var paginationResponse = make(map[string]any)
			if pagination.Start >= (pagination.Limit) {
				prev := fmt.Sprintf("%s%s?start=%d&limit=%d", baseUrl, c.Path(), pagination.Start-pagination.Limit, pagination.Limit)

				if search.Keyword != "" {
					prev += "&keyword=" + search.Keyword
				}

				if sort.Column != "" {
					prev += "&sort=" + sort.Column
				}

				if sort.Direction {
					prev += "&dir=true"
				} else {
					prev += "&dir=false"
				}

				paginationResponse["prev"] = prev
			} else {
				paginationResponse["prev"] = nil
			}

			if totalData > pagination.Start+pagination.Limit {
				next := fmt.Sprintf("%s%s?start=%d&limit=%d", baseUrl, c.Path(), pagination.Start+pagination.Limit, pagination.Limit)

				if search.Keyword != "" {
					next += "&keyword=" + search.Keyword
				}

				if sort.Column != "" {
					next += "&sort=" + sort.Column
				}

				if sort.Direction {
					next += "&dir=true"
				} else {
					next += "&dir=false"
				}

				paginationResponse["next"] = next
			} else {
				paginationResponse["next"] = nil
			}
			response["pagination"] = paginationResponse
		}

		response["message"] = "get all cart success"
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *cartHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var request = new(CartRequest)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		cartId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid cart id"
		}

		userId, err := tokens.ExtractToken(handler.jwtConfig.Secret, token.(*jwt.Token))
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

		var parseInput = new(carts.Cart)
		parseInput.Quantity = request.Quantity

		if err := handler.service.Update(c.Request().Context(), uint(cartId), userId, *parseInput); err != nil {
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

		response["message"] = "update cart success"
		return c.JSON(http.StatusOK, response)
	}
}
