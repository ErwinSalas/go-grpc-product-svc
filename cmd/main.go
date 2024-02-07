package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ErwinSalas/go-grpc-product-svc/pkg/config"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/database"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/product"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/server"
	productpb "github.com/ErwinSalas/go-grpc-product-svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	database := database.Init(c.DBUrl)

	listen, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Product Svc on", c.Port)

	grpcServer := grpc.NewServer()
	productService := product.NewProductService(product.NewGormProductRepository(database.DB)) // Puedes pasar una conexión de base de datos real aquí.
	productpb.RegisterProductServiceServer(grpcServer, server.NewServer(productService))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
