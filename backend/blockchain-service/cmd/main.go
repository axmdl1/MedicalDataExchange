package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/blockchain"
	"github.com/medical-data-exchange/blockchain-service/internal/app"
	"github.com/medical-data-exchange/blockchain-service/internal/config"
	"github.com/medical-data-exchange/blockchain-service/internal/fabric"
	"github.com/medical-data-exchange/blockchain-service/internal/grpc"
	grpcLib "google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()

	// Инициализация Fabric клиента
	fabricClient, err := fabric.NewFabricClient(cfg.Fabric)
	if err != nil {
		log.Fatalf("Failed to create Fabric client: %v", err)
	}
	defer fabricClient.Close()

	// Создаем gRPC сервер
	grpcServer := grpcLib.NewServer()

	// Регистрируем сервис
	blockchainService := app.NewBlockchainService(fabricClient)
	grpcHandler := grpc.NewBlockchainHandler(blockchainService)
	pb.RegisterBlockchainServiceServer(grpcServer, grpcHandler)

	// Запускаем gRPC сервер
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("Blockchain service listening on %s", cfg.GRPCAddr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down blockchain service...")
	grpcServer.GracefulStop()
}
