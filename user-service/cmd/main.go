package main

import (
	"log"

	"github.com/axmdl1/MedicalDataExchange/user-service/internal/app"
)

func main() {
	if err := app.Run("config/local.yaml"); err != nil {
		log.Fatal(err)
	}
}
