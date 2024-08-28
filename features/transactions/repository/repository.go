package repository

import (
	"context"
	"errors"
	"fmt"
	"olshop/features/transactions"
	"olshop/utilities/payment"
	"strconv"
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

	Payment Payment `gorm:"embedded;embeddedPrefix:payment_"`
}

type Payment struct {
	Method        string `gorm:"column:method; type:varchar(20);"`
	Bank          string `gorm:"column:bank; type:varchar(20);"`
	VirtualNumber string `gorm:"column:virtual_number; type:varchar(50);"`
	BillKey       string `gorm:"column:bill_key; type:varchar(50);"`
	BillCode      string `gorm:"column:bill_code; type:varchar(50);"`
	Status        string `gorm:"column:status; type:varchar(20);"`

	CreatedAt time.Time `gorm:"index"`
	ExpiredAt time.Time `gorm:"nullable"`
	PaidAt    time.Time `gorm:"default:null;"`
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

type Cart struct {
	Id       uint  `gorm:"column:id; primaryKey;"`
	Quantity int16 `gorm:"column:quantity; type:int;"`

	ProductId uint    `gorm:"column:product_id"`
	Product   Product `gorm:"foreignKey:ProductId;references:Id"`

	UserId uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:UserId;references:Id"`
}

type transactionRepository struct {
	db      *gorm.DB
	payment payment.Midtrans
}

func NewTransactionRepository(db *gorm.DB, payment payment.Midtrans) transactions.Repository {
	return &transactionRepository{
		db:      db,
		payment: payment,
	}
}

func generateInvoice(userId uint) int {
	now := time.Now()

	invoice := fmt.Sprintf("%d%d%d%d%d%d%d", userId, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	invoiceInt, _ := strconv.Atoi(invoice)
	return invoiceInt
}

func (repo *transactionRepository) Create(ctx context.Context, userId uint, cartIds []uint, newTransaction transactions.Transaction) (*transactions.Transaction, error) {
	var dataCart []Cart
	qry := repo.db.WithContext(ctx).Model(&Cart{}).
		Preload("Product").
		Preload("User").
		Where("carts.user_id = ? AND carts.id IN ?", userId, cartIds)

	if err := qry.Find(&dataCart).Error; err != nil {
		return nil, err
	}

	if len(dataCart) != len(cartIds) {
		return nil, errors.New("one or more cart IDs are invalid")
	}

	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var inputDB = &Transaction{
		Invoice:            generateInvoice(userId),
		PaymentMethod:      newTransaction.PaymentMethod,
		AddressId:          newTransaction.AddressID,
		UserId:             userId,
		Status:             "pending",
		TransactionDetails: []TransactionDetail{},
	}

	for _, cartItem := range dataCart {
		detail := TransactionDetail{
			TransactionInvoice: inputDB.Invoice,
			ProductId:          cartItem.ProductId,
			Quantity:           cartItem.Quantity,
			Subtotal:           float64(cartItem.Quantity) * cartItem.Product.Price,
		}

		inputDB.TransactionDetails = append(inputDB.TransactionDetails, detail)
		inputDB.Total += detail.Subtotal

		deleteQuery := tx.Where("id = ?", cartItem.Id).Delete(&Cart{})
		if deleteQuery.Error != nil {
			tx.Rollback()
			return nil, deleteQuery.Error
		}
	}

	payment, err := repo.payment.NewTransactionPayment(newTransaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	inputDB.Payment = Payment{
		Method:        payment.Method,
		Bank:          payment.Bank,
		VirtualNumber: payment.VirtualNumber,
		BillKey:       payment.BillKey,
		BillCode:      payment.BillCode,
		Status:        payment.Status,
	}

	if err := tx.Create(inputDB).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Mapping to transactions.Transaction type if necessary
	finalTransaction := &transactions.Transaction{
		Invoice:         inputDB.Invoice,
		Total:           inputDB.Total,
		PaymentMethod:   inputDB.PaymentMethod,
		Status:          inputDB.Status,
		TransactionDate: inputDB.TransactionDate,
		UserID:          inputDB.UserId,
		AddressID:       inputDB.AddressId,
		Payment: transactions.Payment{
			Method:        inputDB.Payment.Method,
			Bank:          inputDB.Payment.Bank,
			VirtualNumber: inputDB.Payment.VirtualNumber,
			BillKey:       inputDB.Payment.BillKey,
			BillCode:      inputDB.Payment.BillCode,
			Status:        inputDB.Payment.Status,
			CreatedAt:     inputDB.Payment.CreatedAt,
			ExpiredAt:     inputDB.Payment.ExpiredAt,
			PaidAt:        inputDB.Payment.PaidAt,
		},
	}

	return finalTransaction, nil
}
