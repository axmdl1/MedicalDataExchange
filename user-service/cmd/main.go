package main

import (
	grpchost "github.com/axmdl1/MedicalDataExchange/user-service/internal/app/grpc"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/axmdl1/MedicalDataExchange/user-service/config"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/models"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/repository"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/service"
)

func mustLoad(path string) *config.Config {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}
	var cfg config.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal().Err(err).Msg("parse yaml")
	}
	return &cfg
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cfg := mustLoad("config/local.yaml")

	db, err := gorm.Open(postgres.Open(cfg.Postgres.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("db open")
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal().Err(err).Msg("migrate")
	}

	repo := repository.NewUserRepository(db)
	jwtm := service.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.TTLMin)
	svc := service.NewUserService(repo, jwtm)

	s := grpchost.NewServer(svc)
	if err := s.Start(cfg.GRPC.Port); err != nil {
		log.Fatal().Err(err).Msg("grpc serve")
	}
}
