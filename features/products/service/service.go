package service

import (
	"context"
	"errors"
	"olshop/features/products"
)

type productService struct {
	repo products.Repository
}

func NewProductService(repo products.Repository) products.Service {
	return &productService{
		repo: repo,
	}
}

func (srv *productService) Create(ctx context.Context, data products.Product) error {
	if data.Name == "" {
		return errors.New("validate: name can't be empty")
	}
	if data.Price == 0 {
		return errors.New("validate: price can't be empty")
	}

	if err := srv.repo.Create(ctx, data); err != nil {
		return err
	}

	return nil
}

func (srv *productService) GetAll(ctx context.Context) ([]products.Product, int, error) {
	panic("unimplemented")
}
