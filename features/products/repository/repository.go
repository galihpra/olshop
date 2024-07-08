package repository

import (
	"context"
	"errors"
	"io"
	"olshop/features/products"
	"olshop/helpers/filters"
	"olshop/utilities/cloudinary"

	"gorm.io/gorm"
)

type Product struct {
	Id           uint    `gorm:"column:id; primaryKey;"`
	Name         string  `gorm:"column:name; type:varchar(200);"`
	Price        float64 `gorm:"column:price; type:decimal(16,2);"`
	ThumbnailUrl string  `gorm:"column:thumbnail; type:text;"`
	Rating       float32 `gorm:"column:rating; type:decimal(1,1);"`
	Discount     int     `gorm:"column:discount; type:integer;"`
	Description  string  `gorm:"column:discount; type:text;"`

	Images []Image `gorm:"constraint:OnDelete:CASCADE;"`

	CategoryId uint     `gorm:"column:category_id"`
	Category   Category `gorm:"foreignKey:CategoryId;references:Id"`
}

type Image struct {
	Id        uint      `gorm:"column:id; primaryKey;"`
	ProductId uint      `gorm:"column:product_id;"`
	Product   Product   `gorm:"foreignKey:ProductId;references:Id;constraint:OnDelete:CASCADE;"`
	ImageURL  string    `gorm:"column:image_url; type:text"`
	ImageRaw  io.Reader `gorm:"-"`
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
	inputDB.Discount = data.Discount
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

func (repo *productRepository) GetAll(ctx context.Context, flt filters.Filter) ([]products.Product, int, error) {
	var dataProduct []Product
	var totalData int64

	qry := repo.db.WithContext(ctx).Model(&Product{})

	qry = qry.Select(
		"products.id",
		"products.name",
		"products.price",
		"products.thumbnail",
		"products.discount",
		"products.rating",
	)

	if flt.Search.Keyword != "" {
		qry = qry.Where("name like ?", "%"+flt.Search.Keyword+"%")
	}

	qry.Count(&totalData)

	if flt.Sort.Column != "" {
		dir := "asc"
		if flt.Sort.Direction {
			dir = "desc"
		}

		switch flt.Sort.Column {
		case "rating", "price", "discount":
			qry = qry.Order(flt.Sort.Column + " " + dir)
		default:
			qry = qry.Order("id desc")
		}
	}

	if flt.Pagination.Limit != 0 {
		qry = qry.Limit(flt.Pagination.Limit)
	}

	if flt.Pagination.Start != 0 {
		qry = qry.Offset(flt.Pagination.Start)
	}

	if err := qry.Find(&dataProduct).Error; err != nil {
		return nil, 0, err
	}

	var result []products.Product
	for _, product := range dataProduct {
		result = append(result, products.Product{
			ID:        product.Id,
			Name:      product.Name,
			Price:     product.Price,
			Rating:    product.Rating,
			Discount:  product.Discount,
			Thumbnail: product.ThumbnailUrl,
		})
	}

	return result, int(totalData), nil
}

func (repo *productRepository) Delete(ctx context.Context, id uint) error {
	deleteQuery := repo.db.Delete(&Product{Id: id})
	if deleteQuery.Error != nil {
		return deleteQuery.Error
	}

	if deleteQuery.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func (repo *productRepository) GetProductDetail(ctx context.Context, id uint) (*products.Product, error) {
	var data = new(Product)

	if err := repo.db.Preload("Images").Where("id = ?", id).First(data).Error; err != nil {
		return nil, err
	}

	var result = new(products.Product)
	result.ID = data.Id
	result.Name = data.Name
	result.Price = data.Price
	result.Discount = data.Discount
	result.Description = data.Description

	var images []products.Image
	for _, img := range data.Images {
		images = append(images, products.Image{
			ID:       img.Id,
			ImageURL: img.ImageURL,
		})
	}
	result.Images = images

	return result, nil

}

func (repo *productRepository) Update(ctx context.Context, updateProduct products.Product) error {
	panic("unimplemented")
}
