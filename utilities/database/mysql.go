package database

import (
	"fmt"
	"olshop/config"

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
	// your repositories here
	)

	if err != nil {
		return err
	}

	return nil
}
