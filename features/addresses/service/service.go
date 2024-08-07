package service

import (
	"context"
	"errors"
	"olshop/features/addresses"
)

type addressService struct {
	repo addresses.Repository
}

func NewAddressService(repo addresses.Repository) addresses.Service {
	return &addressService{
		repo: repo,
	}
}

func (service *addressService) Create(ctx context.Context, data addresses.Address) error {
	if data.Street == "" {
		return errors.New("validate: street address can't be empty")
	}
	if data.City == "" {
		return errors.New("validate: city can't be empty")
	}
	if data.Country == "" {
		return errors.New("validate: country can't be empty")
	}
	if data.State == "" {
		return errors.New("validate: state can't be empty")
	}
	if data.Zip == "" {
		return errors.New("validate: zip code can't be empty")
	}

	if err := service.repo.Create(ctx, data); err != nil {
		return err
	}

	return nil
}

func (service *addressService) Delete(ctx context.Context, id uint, userId uint) error {
	if id == 0 {
		return errors.New("validate: invalid id")
	}

	if err := service.repo.Delete(ctx, id, userId); err != nil {
		return err
	}

	return nil
}

func (service *addressService) GetAll(ctx context.Context, userId uint) ([]addresses.Address, error) {
	result, err := service.repo.GetAll(ctx, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
