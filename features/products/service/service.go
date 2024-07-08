package service

import (
	"context"
	"errors"
	"olshop/features/products"
	"olshop/helpers/filters"
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

func (srv *productService) GetAll(ctx context.Context, flt filters.Filter) ([]products.Product, int, error) {
	result, totalData, err := srv.repo.GetAll(ctx, flt)
	if err != nil {
		return nil, 0, err
	}

	return result, totalData, nil
}

func (srv *productService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("validate: invalid id")
	}

	if err := srv.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (srv *productService) GetProductDetail(ctx context.Context, id uint) (*products.Product, error) {
	result, err := srv.repo.GetProductDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (srv *productService) Update(ctx context.Context, updateProduct products.Product) error {
	panic("unimplemented")
}
