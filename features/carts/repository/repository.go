package repository

import (
	"context"
	"errors"
	"olshop/features/carts"
	"olshop/helpers/filters"

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
	Id        uint    `gorm:"column:id; primaryKey;"`
	Name      string  `gorm:"column:name; type:varchar(200);"`
	Price     float64 `gorm:"column:price; type:decimal(16,2);"`
	Thumbnail string  `gorm:"column:thumbnail; type:text;"`
}

type Varian struct {
	Id    uint   `gorm:"column:id; primaryKey;"`
	Color string `gorm:"column:color; type:varchar(20);"`
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

func (repo *cartRepository) GetAll(ctx context.Context, flt filters.Filter, userId uint) ([]carts.Cart, int, error) {
	var dataCart []Cart
	var totalData int64

	qry := repo.db.WithContext(ctx).Model(&Cart{}).
		Preload("Product").
		Preload("Varian").
		Preload("User").
		Where("user_id = ?", userId).
		Order("id DESC")

	if flt.Search.Keyword != "" {
		qry = qry.Where("products.name like ?", "%"+flt.Search.Keyword+"%")
	}

	qry.Count(&totalData)

	if flt.Sort.Column != "" {
		dir := "asc"
		if flt.Sort.Direction {
			dir = "desc"
		}
		qry = qry.Order(flt.Sort.Column + " " + dir)
	}

	if flt.Pagination.Limit != 0 {
		qry = qry.Limit(flt.Pagination.Limit)
	}

	if flt.Pagination.Start != 0 {
		qry = qry.Offset(flt.Pagination.Start)
	}

	if err := qry.Find(&dataCart).Error; err != nil {
		return nil, 0, err
	}

	var result []carts.Cart
	for _, cart := range dataCart {
		result = append(result, carts.Cart{
			ID:       cart.Id,
			Quantity: cart.Quantity,
			Subtotal: cart.Product.Price * float64(cart.Quantity),
			Product: carts.Product{
				ID:        cart.Product.Id,
				Name:      cart.Product.Name,
				Price:     cart.Product.Price,
				Thumbnail: cart.Product.Thumbnail,
			},
			Varian: carts.Varian{
				ID:    cart.Varian.Id,
				Color: cart.Varian.Color,
			},
		})
	}

	return result, int(totalData), nil
}

func (repo *cartRepository) Update(ctx context.Context, cartId uint, userId uint, updateCart carts.Cart) error {
	var model = new(Cart)
	model.Quantity = updateCart.Quantity

	if err := repo.db.Where(&User{Id: userId}).Where(&Cart{Id: cartId}).Updates(model).Error; err != nil {
		return err
	}

	return nil
}
