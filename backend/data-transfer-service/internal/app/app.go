package app

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	grpcapp "github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/app/grpc"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/config"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/repository"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/scheduler"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
)

type App struct {
	GRPC             *grpcapp.App
	CleanupScheduler *scheduler.CleanupScheduler
	ctx              context.Context
	cancel           context.CancelFunc
}

func New(cfg config.Config) *App {
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	repo := repository.NewDataTransferRepository(db)
	svc := service.NewDataTransferService(repo)
	patientAccessService := service.NewPatientAccessService(repo, cfg.Encryption.Secret)
	jwtManager := service.NewJWTManager(cfg.Auth.JWTSecret)

	// Create cleanup scheduler
	cleanupScheduler := scheduler.NewCleanupScheduler(patientAccessService)

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		GRPC:             grpcapp.New(cfg.Server.Port, svc, patientAccessService, jwtManager),
		CleanupScheduler: cleanupScheduler,
		ctx:              ctx,
		cancel:           cancel,
	}
}

func (a *App) Start() {
	// Start cleanup scheduler in background
	go a.CleanupScheduler.Start(a.ctx)

	// Start gRPC server
	go func() {
		if err := a.GRPC.Run(); err != nil {
			panic(err)
		}
	}()
}

func (a *App) Stop() {
	a.cancel() // Stop scheduler
	a.GRPC.Stop()
}
