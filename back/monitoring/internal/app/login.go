package app

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

func (uc *useCase) Login(ctx context.Context, email, password string) (*entities.MonitoringUser, error) {
	user, passwordHash, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domainErrors.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, domainErrors.ErrInvalidCredentials
	}

	if !user.IsPlatformScope() {
		return nil, domainErrors.ErrAccessDenied
	}

	return user, nil
}

func (uc *useCase) GenerateToken(user *entities.MonitoringUser) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
