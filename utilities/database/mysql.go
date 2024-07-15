package database

import (
	"fmt"
	"olshop/config"
	ar "olshop/features/addresses/repository"
	cr "olshop/features/carts/repository"
	pr "olshop/features/products/repository"
	rr "olshop/features/reviews/repository"
	tr "olshop/features/transactions/repository"
	ur "olshop/features/users/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit(cfg config.DatabaseMysql) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MysqlMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&ar.Address{},
		&ur.User{},
		&pr.Category{},
		&pr.Image{},
		&pr.Varian{},
		&pr.Product{},
		&rr.Review{},
		&cr.Cart{},
		&tr.Transaction{},
		&tr.TransactionDetail{},
	)

	if err != nil {
		return err
	}

	return nil
}
