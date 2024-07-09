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

func (service *addressService) Delete(ctx context.Context, id uint) error {
	panic("unimplemented")
}

func (service *addressService) GetAll(ctx context.Context) ([]addresses.Address, error) {
	panic("unimplemented")
}
