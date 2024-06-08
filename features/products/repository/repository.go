package repository

import (
	"context"
	"io"
	"olshop/features/products"
	"olshop/utilities/cloudinary"

	"gorm.io/gorm"
)

type Product struct {
	Id           uint    `gorm:"column:id; primaryKey;"`
	Name         string  `gorm:"column:name; type:varchar(200);"`
	Price        float64 `gorm:"column:price; type:decimal(16,2);"`
	ThumbnailUrl string  `gorm:"column:thumbnail; type:text;"`

	Images []Image `gorm:"many2many:product_images;"`

	CategoryId uint     `gorm:"column:category_id"`
	Category   Category `gorm:"foreignKey:CategoryId;references:Id"`
}

type Image struct {
	Id       uint      `gorm:"column:id; primaryKey;"`
	ImageURL string    `gorm:"column:image_url; type:text"`
	ImageRaw io.Reader `gorm:"-"`
}

type Category struct {
	Id       uint   `gorm:"column:id; primaryKey;"`
	Category string `gorm:"column:category; type:varchar(200);"`
}

type productRepository struct {
	db    *gorm.DB
	cloud cloudinary.Cloud
}

func NewProductRepository(db *gorm.DB, cloud cloudinary.Cloud) products.Repository {
	return &productRepository{
		db:    db,
		cloud: cloud,
	}
}

func (repo *productRepository) Create(ctx context.Context, data products.Product) error {
	var inputDB = new(Product)
	inputDB.Name = data.Name
	inputDB.Price = data.Price
	inputDB.CategoryId = data.Category.ID

	for i := 0; i < len(data.Images); i++ {
		url, err := repo.cloud.Upload(ctx, "products", data.Images[i].ImageRaw)
		if err != nil {
			return err
		}

		image := Image{
			ImageURL: *url,
		}

		switch i {
		case 0:
			inputDB.ThumbnailUrl = image.ImageURL
		}

		inputDB.Images = append(inputDB.Images, image)
	}

	if err := repo.db.Create(inputDB).Error; err != nil {
		return err
	}

	return nil
}

func (repo *productRepository) GetAll(ctx context.Context) ([]products.Product, int, error) {
	panic("unimplemented")
}
