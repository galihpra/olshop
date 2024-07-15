package service

import (
	"context"
	"errors"
	"olshop/features/carts"
)

type cartService struct {
	repo carts.Repository
}

func NewCartService(repo carts.Repository) carts.Service {
	return &cartService{
		repo: repo,
	}
}

func (service *cartService) Create(ctx context.Context, newCart carts.Cart, userId uint) error {
	if newCart.ProductID == 0 {
		return errors.New("validate: product can't be empty")
	}
	if newCart.Quantity == 0 {
		return errors.New("validate: quantity can't be empty")
	}

	if err := service.repo.Create(ctx, newCart, userId); err != nil {
		return err
	}

	return nil
}

func (service *cartService) Delete(ctx context.Context, cartId uint, userId uint) error {
	if cartId == 0 {
		return errors.New("validate: invalid id")
	}

	if err := service.repo.Delete(ctx, cartId, userId); err != nil {
		return err
	}

	return nil
}

func (service *cartService) GetAll(ctx context.Context, userId uint) ([]carts.Cart, error) {
	panic("unimplemented")
}

func (service *cartService) Update(ctx context.Context, cartId uint, userId uint, updateCart carts.Cart) error {
	panic("unimplemented")
}
