package service

import (
	"context"
	"olshop/features/transactions"
)

type transactionService struct {
	repo transactions.Repository
}

func NewTransactionService(repo transactions.Repository) transactions.Service {
	return &transactionService{
		repo: repo,
	}
}

func (service *transactionService) Create(ctx context.Context, userId uint, newTransaction transactions.Transaction) error {
	panic("unimplemented")
}
