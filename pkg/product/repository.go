package product

import (
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	FindOne(id int64) (*models.Product, error)
	DecreaseStock(productID, orderID int64) error
}

type GormProductRepository struct {
	DB *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{
		DB: db,
	}
}

func (r *GormProductRepository) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

func (r *GormProductRepository) FindOne(id int64) (*models.Product, error) {
	var product models.Product
	if result := r.DB.First(&product, id); result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *GormProductRepository) DecreaseStock(productID, orderID int64) error {
	var product models.Product
	if result := r.DB.First(&product, productID); result.Error != nil {
		return result.Error
	}

	if product.Stock <= 0 {
		return nil // You might want to return an error here, depending on your logic
	}

	var log models.StockDecreaseLog
	if result := r.DB.Where(&models.StockDecreaseLog{OrderId: orderID}).First(&log); result.Error == nil {
		return nil // You might want to return an error here, depending on your logic
	}

	product.Stock = product.Stock - 1
	r.DB.Save(&product)

	log.OrderId = orderID
	log.ProductRefer = product.ID
	r.DB.Create(&log)

	return nil
}
