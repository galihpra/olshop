package repository

import (
	"context"
	"errors"
	"io"
	"olshop/features/products"
	"olshop/helpers/filters"
	"olshop/utilities/cloudinary"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id           uint      `gorm:"column:id; primaryKey;"`
	Name         string    `gorm:"column:name; type:varchar(200);"`
	Price        float64   `gorm:"column:price; type:decimal(16,2);"`
	ThumbnailUrl string    `gorm:"column:thumbnail; type:text;"`
	Rating       float32   `gorm:"column:rating; type:float;"`
	Discount     int       `gorm:"column:discount; type:integer;"`
	Description  string    `gorm:"column:description; type:text;"`
	Stock        int       `gorm:"column:stock; type:integer;"`
	DiscountEnd  time.Time `gorm:"column:discount_end; type:timestamp;"`
	Measurement  string    `gorm:"column:measurement; type:varchar(20);"`

	Images []Image `gorm:"constraint:OnDelete:CASCADE;"`

	CategoryId uint     `gorm:"column:category_id"`
	Category   Category `gorm:"foreignKey:CategoryId;references:Id"`
	Varians    []Varian `gorm:"constraint:OnDelete:CASCADE;"`
	Reviews    []Review
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

type Varian struct {
	Id        uint    `gorm:"column:id; primaryKey;"`
	ProductId uint    `gorm:"column:product_id;"`
	Product   Product `gorm:"foreignKey:ProductId;references:Id;constraint:OnDelete:CASCADE;"`
	Color     string  `gorm:"column:color; type:varchar(20);"`
	Stock     int     `gorm:"column:stock; type:integer;"`

	ImageURL string    `gorm:"column:image_url; type:text"`
	ImageRaw io.Reader `gorm:"-"`
}

type Review struct {
	Id        uint      `gorm:"column:id; primaryKey;"`
	ProductId uint      `gorm:"column:product_id;"`
	Product   Product   `gorm:"foreignKey:ProductId;references:Id;constraint:OnDelete:CASCADE;"`
	UserId    uint      `gorm:"column:user_id;"`
	User      User      `gorm:"foreignKey:UserId;references:Id"`
	Review    string    `gorm:"column:review; type:text"`
	Rating    float32   `gorm:"column:rating; type:float"`
	CreatedAt time.Time `gorm:"column:created_at; type:timestamp"`
}

type User struct {
	Id       uint   `gorm:"column:id; primaryKey;"`
	Username string `gorm:"column:username; type:varchar(200)"`
	ImageURL string `gorm:"column:image_url; type:text"`
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
	inputDB.Stock = data.Stock
	inputDB.DiscountEnd = data.DiscountEnd
	inputDB.Measurement = data.Measurement
	inputDB.Description = data.Description

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

	for i := 0; i < len(data.Varians); i++ {
		url, err := repo.cloud.Upload(ctx, "varians", data.Varians[i].ImageRaw)
		if err != nil {
			return err
		}

		varian := Varian{
			Color:    data.Varians[i].Color,
			Stock:    data.Varians[i].Stock,
			ImageURL: *url,
		}

		inputDB.Varians = append(inputDB.Varians, varian)
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

	if err := repo.db.Preload("Images").Preload("Varians").Where("id = ?", id).First(data).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Where("product_id = ?", id).Order("created_at desc").Limit(2).Find(&data.Reviews).Error; err != nil {
		return nil, err
	}

	for i, review := range data.Reviews {
		if err := repo.db.Where("id = ?", review.UserId).First(&data.Reviews[i].User).Error; err != nil {
			return nil, err
		}
	}

	var result = new(products.Product)
	result.ID = data.Id
	result.Name = data.Name
	result.Price = data.Price
	result.Discount = data.Discount
	result.Description = data.Description
	result.Rating = data.Rating
	result.DiscountEnd = data.DiscountEnd
	result.Measurement = data.Measurement
	result.Stock = data.Stock

	var images []products.Image
	for _, img := range data.Images {
		images = append(images, products.Image{
			ID:       img.Id,
			ImageURL: img.ImageURL,
		})
	}
	result.Images = images

	var varians []products.Varian
	for _, varian := range data.Varians {
		varians = append(varians, products.Varian{
			ID:       varian.Id,
			Color:    varian.Color,
			ImageURL: varian.ImageURL,
			Stock:    varian.Stock,
		})
	}
	result.Varians = varians

	var reviews []products.Review
	for _, review := range data.Reviews {
		reviews = append(reviews, products.Review{
			ID:        review.Id,
			Rating:    review.Rating,
			Review:    review.Review,
			CreatedAt: review.CreatedAt,
			User: products.User{
				ID:       review.User.Id,
				Username: review.User.Username,
				ImageURL: review.User.ImageURL,
			},
		})
	}
	result.Reviews = reviews

	return result, nil
}

func (repo *productRepository) Update(ctx context.Context, updateProduct products.Product, id uint) error {
	var existingProduct Product

	if err := repo.db.First(&existingProduct, id).Error; err != nil {
		return err
	}

	if updateProduct.Name != "" {
		existingProduct.Name = updateProduct.Name
	}
	if updateProduct.Price > 0 {
		existingProduct.Price = updateProduct.Price
	}
	if updateProduct.Discount > 0 && updateProduct.Discount <= 100 {
		existingProduct.Discount = updateProduct.Discount
	}
	if updateProduct.Category.ID != 0 {
		existingProduct.CategoryId = updateProduct.Category.ID
	}
	if updateProduct.Stock > 0 {
		existingProduct.Stock = updateProduct.Stock
	}
	if !updateProduct.DiscountEnd.IsZero() {
		existingProduct.DiscountEnd = updateProduct.DiscountEnd
	}
	if updateProduct.Measurement != "" {
		existingProduct.Measurement = updateProduct.Measurement
	}
	if updateProduct.Description != "" {
		existingProduct.Description = updateProduct.Description
	}

	if len(updateProduct.Images) > 0 {
		var images []Image
		for i := 0; i < len(updateProduct.Images); i++ {
			url, err := repo.cloud.Upload(ctx, "products", updateProduct.Images[i].ImageRaw)
			if err != nil {
				return err
			}
			image := Image{
				ImageURL: *url,
			}
			if i == 0 {
				existingProduct.ThumbnailUrl = image.ImageURL
			}
			images = append(images, image)
		}
		existingProduct.Images = images
	}

	if len(updateProduct.Varians) > 0 {
		if err := repo.db.Where("product_id = ?", id).Delete(&Varian{}).Error; err != nil {
			return err
		}

		var variants []Varian
		for i := 0; i < len(updateProduct.Varians); i++ {
			url, err := repo.cloud.Upload(ctx, "varians", updateProduct.Varians[i].ImageRaw)
			if err != nil {
				return err
			}
			variant := Varian{
				Color:    updateProduct.Varians[i].Color,
				Stock:    updateProduct.Varians[i].Stock,
				ImageURL: *url,
			}
			variants = append(variants, variant)
		}
		existingProduct.Varians = variants

	}

	if err := repo.db.Save(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}

func (repo *productRepository) GetAllReview(ctx context.Context, id uint, flt filters.Filter) ([]products.Review, int, error) {
	var dataReview []Review
	var totalData int64

	qry := repo.db.WithContext(ctx).Model(&Review{}).Order("created_at DESC").Where("product_id = ?", id)

	qry.Count(&totalData)

	if flt.Pagination.Limit != 0 {
		qry = qry.Limit(flt.Pagination.Limit)
	}

	if flt.Pagination.Start != 0 {
		qry = qry.Offset(flt.Pagination.Start)
	}

	if err := qry.Find(&dataReview).Error; err != nil {
		return nil, 0, err
	}

	for i, review := range dataReview {
		if err := repo.db.Where("id = ?", review.UserId).First(&dataReview[i].User).Error; err != nil {
			return nil, 0, err
		}
	}

	var result []products.Review
	for _, review := range dataReview {
		result = append(result, products.Review{
			ID:        review.Id,
			Rating:    review.Rating,
			Review:    review.Review,
			CreatedAt: review.CreatedAt,
			User: products.User{
				ID:       review.User.Id,
				Username: review.User.Username,
				ImageURL: review.User.ImageURL,
			},
		})
	}

	return result, int(totalData), nil
}
