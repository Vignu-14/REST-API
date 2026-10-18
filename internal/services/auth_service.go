package services

import (
	"context"
	"errors"
	"rest-api/internal/models"
	"rest-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid username/email or password")

type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}
