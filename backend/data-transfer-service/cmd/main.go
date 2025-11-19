package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/app"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/config"
)

func main() {
	cfg := config.Load()
	application := app.New(cfg)

	// Start application (gRPC server and cleanup scheduler)
	application.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop all services
	application.Stop()
}
