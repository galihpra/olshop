package repository

import (
	"context"
	"errors"
	"olshop/features/carts"

	"gorm.io/gorm"
)

type Cart struct {
	Id       uint  `gorm:"column:id; primaryKey;"`
	Quantity int16 `gorm:"column:quantity; type:int;"`

	ProductId uint    `gorm:"column:product_id"`
	Product   Product `gorm:"foreignKey:ProductId;references:Id"`

	VarianId uint   `gorm:"column:varian_id"`
	Varian   Varian `gorm:"foreignKey:VarianId;references:Id"`

	UserId uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:UserId;references:Id"`
}

type Product struct {
	Id uint `gorm:"column:id; primaryKey;"`
}

type Varian struct {
	Id uint `gorm:"column:id; primaryKey;"`
}

type User struct {
	Id uint `gorm:"column:id; primaryKey;"`
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) carts.Repository {
	return &cartRepository{
		db: db,
	}
}

func (repo *cartRepository) Create(ctx context.Context, newCart carts.Cart, userId uint) error {
	var inputDB = new(Cart)
	inputDB.ProductId = newCart.ProductID
	inputDB.VarianId = newCart.VarianID
	inputDB.Quantity = newCart.Quantity
	inputDB.UserId = userId

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	return nil
}

func (repo *cartRepository) Delete(ctx context.Context, cartId uint, userId uint) error {
	deleteQuery := repo.db.Where(&Cart{UserId: userId}).Delete(&Cart{Id: cartId})
	if deleteQuery.Error != nil {
		return deleteQuery.Error
	}

	if deleteQuery.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func (repo *cartRepository) GetAll(ctx context.Context, userId uint) ([]carts.Cart, error) {
	panic("unimplemented")
}

func (repo *cartRepository) Update(ctx context.Context, cartId uint, userId uint, updateCart carts.Cart) error {
	panic("unimplemented")
}
