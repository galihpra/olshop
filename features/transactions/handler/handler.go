package handler

import (
	"net/http"
	"olshop/config"
	"olshop/features/transactions"
	tokens "olshop/helpers/token"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type transactionHandler struct {
	service   transactions.Service
	jwtConfig config.JWT
}

func NewTransactionHandler(service transactions.Service, jwtConfig config.JWT) transactions.Handler {
	return &transactionHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (handler *transactionHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request = new(TransactionRequest)
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

		var cartIds []uint
		for _, item := range request.Items {
			cartIds = append(cartIds, item.CartID)
		}

		var parseInput = new(transactions.Transaction)
		parseInput.AddressID = request.AddressId
		parseInput.PaymentMethod = request.PaymentMethod

		result, err := handler.service.Create(c.Request().Context(), userId, cartIds, *parseInput)
		if err != nil {
			c.Logger().Error(err)
			if strings.Contains(err.Error(), "validate") {
				response["message"] = strings.ReplaceAll(err.Error(), "validate: ", "")
				return c.JSON(http.StatusBadRequest, response)
			}
			if strings.Contains(err.Error(), "invalid") {
				response["message"] = err.Error()
				return c.JSON(http.StatusBadRequest, response)
			}
			response["message"] = "internal server error"
			return c.JSON(http.StatusInternalServerError, response)
		}

		var data = new(TransactionResponse)
		data.Invoice = result.Invoice
		data.Total = result.Total
		data.Status = result.Status
		data.PaymentBank = result.Payment.Bank
		data.PaymentVirtualNumber = result.Payment.VirtualNumber
		data.PaymentBillKey = result.Payment.BillKey
		data.PaymentBillCode = result.Payment.BillCode
		data.PaymentExpiredAt = &result.Payment.ExpiredAt

		response["data"] = data
		response["message"] = "create transaction success"
		return c.JSON(http.StatusCreated, response)
	}
}
