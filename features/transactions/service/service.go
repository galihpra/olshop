package service

import (
	"context"
	"errors"
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

func (service *transactionService) Create(ctx context.Context, userId uint, cartIds []uint, newTransaction transactions.Transaction) (*transactions.Transaction, error) {
	if newTransaction.PaymentMethod == "" {
		return nil, errors.New("validate: payment method can't be empty")
	}

	if len(cartIds) == 0 {
		return nil, errors.New("validate: no cart selected")
	}

	result, err := service.repo.Create(ctx, userId, cartIds, newTransaction)

	if err != nil {
		return nil, err
	}

	return result, nil
}
