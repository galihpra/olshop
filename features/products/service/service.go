package service

import (
	"context"
	"errors"
	"olshop/features/products"
	"olshop/helpers/filters"
	"time"
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
	if data.Price <= 0 {
		return errors.New("validate: price must be greater than zero")
	}
	if data.Category.ID == 0 {
		return errors.New("validate: category can't be empty")
	}
	if data.Discount <= 0 || data.Discount > 100 {
		return errors.New("validate: discount must be between 1 and 100")
	}
	if data.Discount == 0 {
		return errors.New("validate: discount can't be empty")
	}
	if data.Measurement == "" {
		return errors.New("validate: measurement can't be empty")
	}
	if data.DiscountEnd.Before(time.Now()) {
		return errors.New("validate: discount end must be a future date")
	}

	if len(data.Images) == 0 {
		return errors.New("validate: at least one image is required")
	}

	for _, img := range data.Images {
		if img.ImageRaw == nil {
			return errors.New("validate: all images must have content")
		}
	}

	var totalVarianStock int
	for _, varian := range data.Varians {
		if varian.Color == "" {
			return errors.New("validate: varian color can't be empty")
		}
		if varian.Stock < 0 {
			return errors.New("validate: varian stock can't be negative")
		}
		if varian.ImageRaw == nil {
			return errors.New("validate: varian image must have content")
		}
		totalVarianStock += varian.Stock
	}

	if len(data.Varians) > 0 {
		data.Stock = totalVarianStock
	} else {
		if data.Stock <= 0 {
			return errors.New("validate: stock can't be negative or zero")
		}
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

func (srv *productService) Update(ctx context.Context, updateProduct products.Product, id uint) error {
	if id == 0 {
		return errors.New("validate: invalid product id")
	}
	if updateProduct.Price < 0 {
		return errors.New("validate: price must be greater than zero")
	}
	if updateProduct.Discount < 0 || updateProduct.Discount > 100 {
		return errors.New("validate: discount must be between 1 and 100")
	}
	if updateProduct.Discount != 0 {
		if updateProduct.DiscountEnd.Before(time.Now()) {
			return errors.New("validate: discount end must be a future date")
		}
	}

	var totalVarianStock int
	for _, varian := range updateProduct.Varians {
		if varian.Stock < 0 {
			return errors.New("validate: varian stock can't be negative")
		}
		totalVarianStock += varian.Stock
	}

	if len(updateProduct.Varians) > 0 {
		updateProduct.Stock = totalVarianStock
	} else {
		if updateProduct.Stock < 0 {
			return errors.New("validate: stock can't be negative or zero")
		}
	}

	if err := srv.repo.Update(ctx, updateProduct, id); err != nil {
		return err
	}

	return nil
}

func (srv *productService) GetAllReview(ctx context.Context, id uint, flt filters.Filter) ([]products.Review, int, error) {
	result, totalData, err := srv.repo.GetAllReview(ctx, id, flt)
	if err != nil {
		return nil, 0, err
	}

	return result, totalData, nil
}
