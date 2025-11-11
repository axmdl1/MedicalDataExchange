package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/axmdl1/MedicalDataExchange/user-service/internal/models"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	Get(ctx context.Context, id int64) (*models.User, error)
	List(ctx context.Context, clinicID *int64, userType string) ([]models.User, error)
	Update(ctx context.Context, u *models.User) (*models.User, error)
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (*models.User, string, int64, error)
}

type userService struct {
	repo repository.UserRepository
	jwt  *JWTManager
}

func NewUserService(repo repository.UserRepository, jwt *JWTManager) UserService {
	return &userService{repo: repo, jwt: jwt}
}

func (s *userService) Create(ctx context.Context, u *models.User) (*models.User, error) {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	// Если есть soft-deleted аккаунт — просим зайти (Login) для восстановления
	if _, err := s.repo.GetSoftDeletedByEmail(ctx, u.Email); err == nil {
		return nil, errors.New("account for this email was deleted earlier; please login to restore it")
	}

	if u.Type == "employee" && u.ClinicID == nil {
		return nil, errors.New("clinic_id is required for employee")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)
	if err := s.repo.Create(ctx, u); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("email already used")
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) Get(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context, clinicID *int64, userType string) ([]models.User, error) {
	if clinicID != nil {
		return s.repo.List(ctx, *clinicID, true, userType)
	}
	return s.repo.List(ctx, 0, false, userType)
}

func (s *userService) Update(ctx context.Context, u *models.User) (*models.User, error) {
	if u.Type == "employee" && u.ClinicID == nil {
		return nil, errors.New("clinic_id is required for employee")
	}
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(hash)
	}
	return u, s.repo.Update(ctx, u)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Login(ctx context.Context, email, password string) (*models.User, string, int64, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// 1) ищем активного
	u, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
			return nil, "", 0, ErrInvalidCredentials
		}
		token, exp, err := s.jwt.Sign(u.ID, u.Type, u.ClinicID, u.Email)
		if err != nil {
			return nil, "", 0, err
		}
		return u, token, exp.Unix(), nil
	}

	// 2) не нашли активного — пробуем soft-deleted
	deleted, err2 := s.repo.GetSoftDeletedByEmail(ctx, email)
	if err2 != nil {
		// вообще нет такого e-mail
		return nil, "", 0, ErrInvalidCredentials
	}
	// проверяем пароль от старого аккаунта
	if bcrypt.CompareHashAndPassword([]byte(deleted.Password), []byte(password)) != nil {
		return nil, "", 0, ErrInvalidCredentials
	}
	// реактивируем
	if err := s.repo.Reactivate(ctx, deleted.ID); err != nil {
		return nil, "", 0, err
	}
	token, exp, err := s.jwt.Sign(deleted.ID, deleted.Type, deleted.ClinicID, deleted.Email)
	if err != nil {
		return nil, "", 0, err
	}
	return deleted, token, exp.Unix(), nil
}
