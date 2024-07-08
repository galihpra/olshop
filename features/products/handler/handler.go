package handler

import (
	"context"
	"fmt"
	"net/http"
	"olshop/features/products"
	"olshop/helpers/filters"
	"strings"

	"github.com/labstack/echo/v4"
)

type productHandler struct {
	service products.Service
}

func NewProductHandler(service products.Service) products.Handler {
	return &productHandler{
		service: service,
	}
}

func (hdl *productHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(CreateRequest)
		var response = make(map[string]any)

		token := c.Get("user")
		if token == nil {
			response["message"] = "unauthorized access"
			return c.JSON(http.StatusUnauthorized, response)
		}

		if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
			c.Logger().Error(err)
			response["message"] = "failed to parse form data"
			return c.JSON(http.StatusBadRequest, response)
		}

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "incorect input data"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(products.Product)
		parseInput.Name = request.Name
		parseInput.Price = request.Price
		parseInput.Discount = request.Discount
		parseInput.Category.ID = request.CategoryId

		// Handle file uploads
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["images"]
			for _, file := range files {
				src, err := file.Open()
				if err != nil {
					c.Logger().Error(err)
					response["message"] = "failed to open file"
					return c.JSON(http.StatusInternalServerError, response)
				}
				defer src.Close()

				request.Images = append(request.Images, src)
			}
		}

		for _, file := range request.Images {
			if file != nil {
				parseInput.Images = append(parseInput.Images, products.Image{
					ImageRaw: file,
				})
			}
		}

		if err := hdl.service.Create(c.Request().Context(), *parseInput); err != nil {
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

		response["message"] = "create product success"
		return c.JSON(http.StatusCreated, response)
	}
}

func (hdl *productHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var baseUrl = c.Scheme() + "://" + c.Request().Host

		var pagination = new(filters.Pagination)
		c.Bind(pagination)
		if pagination.Start != 0 && pagination.Limit == 0 {
			pagination.Limit = 5
		}

		var search = new(filters.Search)
		c.Bind(search)

		var sort = new(filters.Sort)
		c.Bind(sort)

		result, totalData, err := hdl.service.GetAll(context.Background(), filters.Filter{Search: *search, Pagination: *pagination, Sort: *sort})
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data []ProductResponse
		for _, product := range result {
			data = append(data, ProductResponse{
				Id:        product.ID,
				Name:      product.Name,
				Price:     product.Price,
				Rating:    product.Rating,
				Discount:  product.Discount,
				Thumbnail: product.Thumbnail,
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

		response["message"] = "get all product success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *productHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic("unimplemented")
	}
}

func (hdl *productHandler) GetProductDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic("unimplemented")
	}
}

func (hdl *productHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic("unimplemented")
	}
}
