package server

import (
	"context"

	"github.com/ErwinSalas/go-grpc-product-svc/pkg/product"
	productpb "github.com/ErwinSalas/go-grpc-product-svc/pkg/proto"
)

type Server struct {
	productpb.UnimplementedProductServiceServer
	Service *product.ProductService
}

func NewServer(service *product.ProductService) *Server {
	return &Server{
		Service: service,
	}
}

func (s *Server) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	return s.Service.CreateProduct(req.Name, req.Stock, req.Price)
}

func (s *Server) FindOne(ctx context.Context, req *productpb.FindOneRequest) (*productpb.FindOneResponse, error) {
	return s.Service.FindOneProductByID(req.Id)
}

func (s *Server) DecreaseStock(ctx context.Context, req *productpb.DecreaseStockRequest) (*productpb.DecreaseStockResponse, error) {
	return s.Service.DecreaseStock(req.Id, req.OrderId)
}
