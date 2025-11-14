package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	userv1 "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/user"
	appcfg "github.com/axmdl1/MedicalDataExchange/user-service/config"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/app/grpc"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/models"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/repository"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/service"
)

func loadConfig(path string) (*appcfg.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg appcfg.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// пул соединений
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	// быстрая проверка коннекта
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func Run(configPath string) error {
	// логгер
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// конфиг
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	log.Info().Str("port", cfg.GRPC.Port).Msg("starting user-service")

	// база
	db, err := openDB(cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	// миграции
	if err := db.AutoMigrate(&models.User{}); err != nil {
		// GORM может попытаться DROP CONSTRAINT "uni_users_email"
		// Если его нет (SQLSTATE 42704) — просто проигнорируем.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42704" && strings.Contains(strings.ToLower(pgErr.Message), "uni_users_email") {
			log.Warn().Err(err).Msg("ignoring missing constraint uni_users_email during migration")
		} else {
			return err
		}
	}

	// гарантируем глобальную уникальность email (регистр-независимо)
	if err := db.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_idx
    ON users (lower(email));`).Error; err != nil {
		return err
	}

	// зависимости
	userRepo := repository.NewUserRepository(db)
	jwt := service.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.TTLMin)
	userSvc := service.NewUserService(userRepo, jwt)

	// gRPC сервер
	s := grpc.NewServer(userSvc, jwt)
	log.Info().Str("addr", cfg.GRPC.Port).Msg("grpc listen")
	return s.Start(cfg.GRPC.Port)
}

// Ensure imports are used (so go vet is happy)
var _ = userv1.User{}
