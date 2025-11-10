package service

import (
	"context"
	"errors"

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
	List(ctx context.Context, clinicID int64, userType string) ([]models.User, error)
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
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Get(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context, clinicID int64, userType string) ([]models.User, error) {
	return s.repo.List(ctx, clinicID, userType)
}

func (s *userService) Update(ctx context.Context, u *models.User) (*models.User, error) {
	// если пришёл новый пароль — захешируем
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
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", 0, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, "", 0, ErrInvalidCredentials
	}
	token, exp, err := s.jwt.Sign(u.ID, u.Type)
	if err != nil {
		return nil, "", 0, err
	}
	return u, token, exp.Unix(), nil
}
