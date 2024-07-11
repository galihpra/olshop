package repository

import (
	"olshop/features/reviews"
	"time"

	"gorm.io/gorm"
)

type Review struct {
	Id     uint    `gorm:"column:id; primaryKey;"`
	Review string  `gorm:"column:review; type:text;"`
	Rating float32 `gorm:"column:rating; type:float;"`

	ProductId uint    `gorm:"column:product_id"`
	Product   Product `gorm:"foreignKey:ProductId;references:Id"`

	UserId uint `gorm:"column:user_id"`
	User   User `gorm:"foreignKey:UserId;references:Id"`

	CreatedAt time.Time `gorm:"column:created_at; type:timestamp;"`
}

type User struct {
	Id       uint     `gorm:"column:id; primaryKey;"`
	Username string   `gorm:"column:username;"`
	Image    string   `gorm:"column:image;"`
	Reviews  []Review `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
}

type Product struct {
	Id      uint     `gorm:"column:id; primaryKey;"`
	Rating  float32  `gorm:"column:rating; type:float;"`
	Reviews []Review `gorm:"foreignKey:ProductId;references:Id;constraint:OnDelete:CASCADE;"`
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) reviews.Repository {
	return &reviewRepository{
		db: db,
	}
}

func (repo *reviewRepository) Create(userId uint, newReview reviews.Review) error {
	var inputDB = new(Review)
	inputDB.Review = newReview.Review
	inputDB.ProductId = newReview.ProductId
	inputDB.Rating = newReview.Rating
	inputDB.UserId = userId

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	if err := repo.UpdateProductRating(newReview.ProductId); err != nil {
		return err
	}

	return nil
}

func (repo *reviewRepository) UpdateProductRating(productId uint) error {
	var product Product
	if err := repo.db.Preload("Reviews").First(&product, productId).Error; err != nil {
		return err
	}

	var totalRating float32
	for _, review := range product.Reviews {
		totalRating += review.Rating
	}

	product.Rating = totalRating / float32(len(product.Reviews))

	if err := repo.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
