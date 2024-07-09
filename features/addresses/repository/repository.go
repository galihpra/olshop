package repository

import (
	"context"
	"olshop/features/addresses"

	"gorm.io/gorm"
)

type Address struct {
	Id      uint   `gorm:"column:id; primaryKey;"`
	Street  string `gorm:"column:street; type:varchar(200);"`
	Country string `gorm:"column:country; type:varchar(30);"`
	City    string `gorm:"column:city; type:varchar(45);"`
	State   string `gorm:"column:state; type:varchar(45);"`
	Zip     string `gorm:"column:zip_code; type:varchar(10);"`

	UserId uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:UserId;references:Id"`
}

type User struct {
	Id        uint      `gorm:"column:id; primaryKey;"`
	Adrresses []Address `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) addresses.Repository {
	return &addressRepository{
		db: db,
	}
}

func (repo *addressRepository) Create(ctx context.Context, data addresses.Address) error {
	var inputDB = new(Address)
	inputDB.Street = data.Street
	inputDB.City = data.City
	inputDB.Country = data.Country
	inputDB.State = data.State
	inputDB.Zip = data.Zip
	inputDB.UserId = data.UserID

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	return nil
}

func (repo *addressRepository) Delete(ctx context.Context, id uint) error {
	panic("unimplemented")
}

func (repo *addressRepository) GetAll(ctx context.Context) ([]addresses.Address, error) {
	panic("unimplemented")
}
