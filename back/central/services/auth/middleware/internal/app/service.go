package app

import (
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/domain"
)

type AuthService struct {
	jwtService domain.IJWTService
}

func NewAuthService(jwtService domain.IJWTService) *AuthService {
	return &AuthService{jwtService: jwtService}
}

// ValidateBusinessToken validates a token and ensures it is a business token
func (s *AuthService) ValidateBusinessToken(token string) (*domain.AuthInfo, error) {
	if token == "" {
		return nil, &domain.AuthError{Message: "Token de autorización requerido"}
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Try to validate as business token
	businessClaims, err := s.jwtService.ValidateBusinessToken(token)
	if err == nil {
		if businessClaims.TokenType != "business" {
			return nil, &domain.AuthError{Message: fmt.Sprintf("Token type inválido: %s, se requiere 'business'", businessClaims.TokenType)}
		}
	} else {
		// If not business token, try to validate as main token to give better error
		mainClaims, err2 := s.jwtService.ValidateToken(token)
		if err2 == nil {
			if mainClaims.TokenType == "main" {
				return nil, &domain.AuthError{Message: "Este endpoint requiere un business token, no un token principal"}
			}
		}
		return nil, &domain.AuthError{Message: "Se requiere un business token válido"}
	}

	return &domain.AuthInfo{
		Type:                domain.AuthTypeJWT,
		UserID:              businessClaims.UserID,
		BusinessID:          businessClaims.BusinessID,
		BusinessTokenClaims: businessClaims,
	}, nil
}

// ValidateMainToken validates a token and ensures it is a main token (for BusinessTokenAuth)
func (s *AuthService) ValidateMainToken(token string) (*domain.AuthInfo, error) {
	if token == "" {
		return nil, &domain.AuthError{Message: "Token de autorización requerido"}
	}

	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	mainClaims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, &domain.AuthError{Message: "Token inválido"}
	}

	if mainClaims.TokenType != "main" {
		return nil, &domain.AuthError{Message: "Se requiere un token principal, no un business token"}
	}

	return &domain.AuthInfo{
		Type:      domain.AuthTypeJWT,
		UserID:    mainClaims.UserID,
		JWTClaims: mainClaims,
	}, nil
}
