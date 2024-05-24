package repository

import (
	"olshop/features/users"
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
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.Repository {
	return &userRepository{
		db: db,
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
	result.Username = data.Username
	result.Email = data.Email

	return result, nil
}
