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
	app := app.New(cfg)

	// Запуск
	go func() {
		if err := app.GRPC.Run(); err != nil {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.GRPC.Stop()
}
