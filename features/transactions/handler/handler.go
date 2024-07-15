package handler

import (
	"olshop/config"
	"olshop/features/transactions"

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
	panic("unimplemented")
}
