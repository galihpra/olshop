package repository

import (
	"context"
	"olshop/features/transactions"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	Invoice         int       `gorm:"column:invoice;primaryKey"`
	Total           float64   `gorm:"column:total;type:decimal(16,2)"`
	PaymentMethod   string    `gorm:"column:payment_method;type:varchar(20)"`
	Status          string    `gorm:"column:status;type:varchar(20)"`
	TransactionDate time.Time `gorm:"column:transaction_date;type:timestamp;default:current_timestamp()"`

	UserId uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:UserId;references:Id"`

	AddressId uint    `gorm:"column:address_id"`
	Address   Address `gorm:"foreignKey:AddressId;references:Id"`

	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionInvoice;constraint:OnDelete:CASCADE;"`
}

type User struct {
	Id uint `gorm:"column:id;primaryKey"`
}

type TransactionDetail struct {
	Id       uint    `gorm:"column:id;primaryKey"`
	Quantity int16   `gorm:"column:quantity;type:integer"`
	Subtotal float64 `gorm:"column:sub_total;type:decimal(16,2)"`

	TransactionInvoice int         `gorm:"column:transaction_invoice"`
	Transaction        Transaction `gorm:"foreignKey:TransactionInvoice;references:Invoice"`

	ProductId uint    `gorm:"column:product_id"`
	Product   Product `gorm:"foreignKey:ProductId;references:Id"`
}

type Address struct {
	Id uint `gorm:"column:id;primaryKey"`
}

type Product struct {
	Id    uint    `gorm:"column:id;primaryKey"`
	Name  string  `gorm:"column:name;type:varchar(200)"`
	Price float64 `gorm:"column:price;type:decimal(16,2)"`
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) transactions.Repository {
	return &transactionRepository{
		db: db,
	}
}

func (repo *transactionRepository) Create(ctx context.Context, userId uint, newTransaction transactions.Transaction) error {
	panic("unimplemented")
}
