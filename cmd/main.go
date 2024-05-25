package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ErwinSalas/go-grpc-product-svc/pkg/config"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/database"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/product"
	"github.com/ErwinSalas/go-grpc-product-svc/pkg/server"
	productpb "github.com/ErwinSalas/go-grpc-product-svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request: %v", req)
	log.Printf("Method: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	return resp, err
}

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
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)

	// Start health check routine
	go func() {
		for {
			var count int64
			if err := database.DB.Table("products").Count(&count).Error; err != nil {
				log.Println("Database query error:", err)
				healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
				return
			} else {
				healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
				time.Sleep(5 * time.Second)
			}
		}
	}()
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
