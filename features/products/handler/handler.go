package handler

import (
	"context"
	"fmt"
	"net/http"
	"olshop/features/products"
	"olshop/helpers/filters"
	"strconv"
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
			response["message"] = "incorrect input data"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(products.Product)
		parseInput.Name = request.Name
		parseInput.Price = request.Price
		parseInput.Discount = request.Discount
		parseInput.Description = request.Description
		parseInput.Category.ID = request.CategoryId
		parseInput.Stock = request.Stock
		parseInput.Measurement = request.Measurement
		parseInput.DiscountEnd = request.DiscountEnd

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

				parseInput.Images = append(parseInput.Images, products.Image{
					ImageRaw: src,
				})
			}

			for i := 0; i < len(form.File); i++ {
				colorKey := fmt.Sprintf("varians[%d].color", i)
				stockKey := fmt.Sprintf("varians[%d].stock", i)
				varianImageKey := fmt.Sprintf("varians[%d].varian_image", i)

				color := c.FormValue(colorKey)
				stock, _ := strconv.Atoi(c.FormValue(stockKey))
				file, _ := c.FormFile(varianImageKey)
				if file != nil {
					src, err := file.Open()
					if err != nil {
						c.Logger().Error(err)
						response["message"] = "failed to open varian image file"
						return c.JSON(http.StatusInternalServerError, response)
					}
					defer src.Close()

					parseInput.Varians = append(parseInput.Varians, products.Varian{
						Color:    color,
						Stock:    stock,
						ImageRaw: src,
					})
				}
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
		var response = make(map[string]any)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid product id"
		}

		if err := hdl.service.Delete(c.Request().Context(), uint(id)); err != nil {
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

		response["message"] = "delete product success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *productHandler) GetProductDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			response["message"] = "invalid product id"
			return c.JSON(http.StatusBadRequest, response)
		}

		result, err := hdl.service.GetProductDetail(c.Request().Context(), uint(id))
		if err != nil {
			c.Logger().Error(err)

			if strings.Contains(err.Error(), "not found") {
				response["message"] = "not found"
				return c.JSON(http.StatusNotFound, response)
			}

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data = new(ProductResponse)
		data.Id = result.ID
		data.Name = result.Name
		data.Price = result.Price
		data.Description = result.Description
		data.Discount = result.Discount
		data.Rating = result.Rating
		data.Stock = result.Stock
		data.Measurement = result.Measurement
		data.DiscountEnd = result.DiscountEnd

		var images []string
		for _, img := range result.Images {
			images = append(images, img.ImageURL)
		}
		data.Images = images

		var varians []Varianresponse
		for _, varian := range result.Varians {
			varians = append(varians, Varianresponse{
				Id:       varian.ID,
				Color:    varian.Color,
				Stock:    varian.Stock,
				ImageURL: varian.ImageURL,
			})
		}
		data.Varians = varians

		var reviews []ReviewResponse
		for _, review := range result.Reviews {
			reviews = append(reviews, ReviewResponse{
				ID:        review.ID,
				Review:    review.Review,
				Rating:    review.Rating,
				CreatedAt: review.CreatedAt,
				User: UserResponse{
					UserID:   review.User.ID,
					Username: review.User.Username,
					ImageURL: review.User.ImageURL,
				},
			})
		}
		data.Reviews = reviews

		response["data"] = data
		response["message"] = "get detail product success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *productHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var request = new(CreateRequest)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "invalid product id"
		}

		if err := c.Bind(request); err != nil {
			c.Logger().Error(err)

			response["message"] = "please fill input correctly"
			return c.JSON(http.StatusBadRequest, response)
		}

		var parseInput = new(products.Product)
		parseInput.Name = request.Name
		parseInput.Price = request.Price
		parseInput.Discount = request.Discount
		parseInput.Description = request.Description
		parseInput.Category.ID = request.CategoryId
		parseInput.Stock = request.Stock
		parseInput.Measurement = request.Measurement
		parseInput.DiscountEnd = request.DiscountEnd

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

				parseInput.Images = append(parseInput.Images, products.Image{
					ImageRaw: src,
				})
			}

			for i := 0; i < len(form.File); i++ {
				colorKey := fmt.Sprintf("varians[%d].color", i)
				stockKey := fmt.Sprintf("varians[%d].stock", i)
				varianImageKey := fmt.Sprintf("varians[%d].varian_image", i)

				color := c.FormValue(colorKey)
				stock, _ := strconv.Atoi(c.FormValue(stockKey))
				file, _ := c.FormFile(varianImageKey)
				if file != nil {
					src, err := file.Open()
					if err != nil {
						c.Logger().Error(err)
						response["message"] = "failed to open varian image file"
						return c.JSON(http.StatusInternalServerError, response)
					}
					defer src.Close()

					parseInput.Varians = append(parseInput.Varians, products.Varian{
						Color:    color,
						Stock:    stock,
						ImageRaw: src,
					})
				}
			}
		}

		if err := hdl.service.Update(c.Request().Context(), *parseInput, uint(id)); err != nil {
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

		response["message"] = "update product success"
		return c.JSON(http.StatusOK, response)
	}
}

func (hdl *productHandler) GetAllReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response = make(map[string]any)
		var baseUrl = c.Scheme() + "://" + c.Request().Host

		var pagination = new(filters.Pagination)
		c.Bind(pagination)
		if pagination.Start != 0 && pagination.Limit == 0 {
			pagination.Limit = 5
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			response["message"] = "invalid product id"
			return c.JSON(http.StatusBadRequest, response)
		}

		result, totalData, err := hdl.service.GetAllReview(context.Background(), uint(id), filters.Filter{Pagination: *pagination})
		if err != nil {
			c.Logger().Error(err)

			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var reviews []ReviewResponse
		for _, review := range result {
			reviews = append(reviews, ReviewResponse{
				ID:        review.ID,
				Review:    review.Review,
				Rating:    review.Rating,
				CreatedAt: review.CreatedAt,
				User: UserResponse{
					UserID:   review.User.ID,
					Username: review.User.Username,
					ImageURL: review.User.ImageURL,
				},
			})
		}
		response["data"] = reviews

		if pagination.Limit != 0 {
			var paginationResponse = make(map[string]any)
			if pagination.Start >= pagination.Limit {
				prev := fmt.Sprintf("%s%s?start=%d&limit=%d", baseUrl, c.Path(), pagination.Start-pagination.Limit, pagination.Limit)
				paginationResponse["prev"] = strings.Replace(prev, ":id", strconv.Itoa(id), 1)
			} else {
				paginationResponse["prev"] = nil
			}

			if totalData > pagination.Start+pagination.Limit {
				next := fmt.Sprintf("%s%s?start=%d&limit=%d", baseUrl, c.Path(), pagination.Start+pagination.Limit, pagination.Limit)
				paginationResponse["next"] = strings.Replace(next, ":id", strconv.Itoa(id), 1)
			} else {
				paginationResponse["next"] = nil
			}
			response["pagination"] = paginationResponse
		}

		response["message"] = "get all review success"
		return c.JSON(http.StatusOK, response)
	}
}
