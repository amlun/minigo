package service

import (
	"context"

	"minigo/internal/domain/entity"
	"minigo/internal/domain/repository"
	"minigo/internal/infrastructure/auth"
	"minigo/internal/infrastructure/config"
)

// AuthService provides authentication operations.
type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{userRepo: users}
}

// Login validates credentials and returns JWT token and user info.
func (s *AuthService) Login(ctx context.Context, phone, password string) (string, error) {
	var (
		err  error
		user *entity.User
	)
	if user, err = s.userRepo.GetByPhone(ctx, phone); err != nil {
		return "", ErrUserNotFound
	}
	// Verify password
	if err = user.CheckPassword(password); err != nil {
		return "", ErrInvalidCredentials
	}
	// Generate JWT with correct user role. TODO
	userRole := entity.RoleUser
	token, err := auth.GenerateToken(user.ID, userRole, config.GetJWTExpireDuration())
	if err != nil {
		return "", err
	}
	return token, nil
}
