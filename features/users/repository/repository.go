package repository

import (
	"context"
	"olshop/features/users"
	"olshop/utilities/cloudinary"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        uint   `gorm:"column:id; primaryKey;"`
	Name      string `gorm:"column:name; type:varchar(200);"`
	Email     string `gorm:"column:email; type:varchar(20);unique"`
	Password  string `gorm:"column:password; type:varchar(72);"`
	Image     string `gorm:"column:image; type:text;"`
	Username  string `gorm:"column:username; type:varchar(45);unique"`
	CreatedAt time.Time
}

type userRepository struct {
	db    *gorm.DB
	cloud cloudinary.Cloud
}

func NewUserRepository(db *gorm.DB, cloud cloudinary.Cloud) users.Repository {
	return &userRepository{
		db:    db,
		cloud: cloud,
	}
}

func (repo *userRepository) Register(newUser users.User) error {
	var inputDB = new(User)
	inputDB.Name = newUser.Name
	inputDB.Email = newUser.Email
	inputDB.Password = newUser.Password
	inputDB.Username = newUser.Username

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Login(email string) (*users.User, error) {
	var data = new(User)

	if err := repo.db.Where("email = ?", email).First(data).Error; err != nil {
		return nil, err
	}

	var result = new(users.User)
	result.Id = data.Id
	result.Name = data.Name
	result.Password = data.Password
	result.Username = data.Username
	result.Email = data.Email

	return result, nil
}

func (repo *userRepository) Update(id uint, updateUser users.User) error {
	var model = new(User)
	model.Name = updateUser.Name
	model.Email = updateUser.Email
	model.Password = updateUser.Password
	model.Username = updateUser.Username

	url, err := repo.cloud.Upload(context.Background(), "users", updateUser.ImageRaw)
	if err != nil {
		return err
	}
	model.Image = *url

	if err := repo.db.Where(&User{Id: id}).Updates(model).Error; err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Delete(id uint) error {
	panic("unimplemented")
}

func (repo *userRepository) GetById(id uint) (*users.User, error) {
	panic("unimplemented")
}
