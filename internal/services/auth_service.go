package services

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/models"
	"rest-api/internal/repository"
	"rest-api/pkg/password"
	"rest-api/pkg/token"
	"strings"
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

func (s *authService) Login (ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Identifier == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// Fetch user by username or email
	user, err := s.userRepo.GetByIdentifier(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to fetch user login: %w", err)
	}

	// Comapre the password
	if err := password.Compare(req.Password, user.PasswordHash); err != nil{
		return nil, ErrInvalidCredentials
	}

	jwtToken, err := token.GenerateToken(user.ID, string(user.Role), s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: jwtToken,
		Role: user.Role,
		ID: user.ID,
	}, nil

}
