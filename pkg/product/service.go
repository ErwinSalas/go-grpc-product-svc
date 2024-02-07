package product

import (
	"net/http"

	"github.com/ErwinSalas/go-grpc-product-svc/pkg/models"
	productpb "github.com/ErwinSalas/go-grpc-product-svc/pkg/proto"
)

type ProductService struct {
	Repository ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		Repository: repo,
	}
}

func (s *ProductService) CreateProduct(name string, stock, price int64) (*productpb.CreateProductResponse, error) {
	var product models.Product
	product.Name = name
	product.Stock = stock
	product.Price = price

	if err := s.Repository.CreateProduct(&product); err != nil {
		return &productpb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}, nil
	}

	return &productpb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.ID,
	}, nil
}

func (s *ProductService) FindOneProductByID(id int64) (*productpb.FindOneResponse, error) {
	product, err := s.Repository.FindOne(id)
	if err != nil {
		return &productpb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		}, nil
	}

	data := &productpb.FindOneData{
		Id:    product.StockDecreaseLogs.ID,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &productpb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *ProductService) DecreaseStock(id, orderId int64) (*productpb.DecreaseStockResponse, error) {
	if err := s.Repository.DecreaseStock(id, orderId); err != nil {
		return &productpb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		}, nil
	}

	return &productpb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
