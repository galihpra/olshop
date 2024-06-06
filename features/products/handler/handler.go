package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/products"
	"strings"

	"github.com/labstack/echo/v4"
)

type productHandler struct {
	service   products.Service
	jwtConfig config.JWT
}

func NewProductHandler(service products.Service, jwtConfig config.JWT) products.Handler {
	return &productHandler{
		service:   service,
		jwtConfig: jwtConfig,
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
		panic("unimplemented")
	}
}
