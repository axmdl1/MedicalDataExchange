package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/app"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/config"
)

func main() {
	cfg := config.Load()
	app := app.New(cfg)

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("\nshutting down...")
}
