package app

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	grpcapp "github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/app/grpc"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/config"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/repository"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
)

type App struct {
	GRPC *grpcapp.App
}

func New(cfg config.Config) *App {
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	repo := repository.NewDataTransferRepository(db)
	svc := service.NewDataTransferService(repo)
	jwtManager := service.NewJWTManager(cfg.Auth.JWTSecret)

	return &App{
		GRPC: grpcapp.New(cfg.Server.Port, svc, jwtManager),
	}
}
