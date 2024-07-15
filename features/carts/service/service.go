package service

import (
	"context"
	"errors"
	"olshop/features/carts"
	"olshop/helpers/filters"
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
	if newCart.Quantity <= 0 {
		return errors.New("validate: quantity can't be zero or minus")
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

func (service *cartService) GetAll(ctx context.Context, flt filters.Filter, userId uint) ([]carts.Cart, int, error) {
	if userId == 0 {
		return nil, 0, errors.New("validate: invalid id")
	}

	result, totalData, err := service.repo.GetAll(ctx, flt, userId)
	if err != nil {
		return nil, 0, err
	}

	return result, totalData, nil
}

func (service *cartService) Update(ctx context.Context, cartId uint, userId uint, updateCart carts.Cart) error {
	if cartId == 0 {
		return errors.New("validate: invalid cart id")
	}
	if userId == 0 {
		return errors.New("validate: invalid user id")
	}
	if updateCart.Quantity <= 0 {
		return errors.New("validate: quantity can't be zero or minus")
	}

	if err := service.repo.Update(ctx, cartId, userId, updateCart); err != nil {
		return err
	}

	return nil
}
