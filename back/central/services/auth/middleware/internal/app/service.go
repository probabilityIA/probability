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

// ValidateToken validates the unified token
func (s *AuthService) ValidateToken(token string) (*domain.AuthInfo, error) {
	if token == "" {
		return nil, &domain.AuthError{Message: "Token de autorización requerido"}
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	// Validate unified token
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, &domain.AuthError{Message: fmt.Sprintf("Token inválido: %v", err)}
	}

	// Debug log
	fmt.Printf("[DEBUG JWT] UserID: %d, BusinessID: %d, RoleID: %d\n", claims.UserID, claims.BusinessID, claims.RoleID)

	businessTokenClaims := &domain.BusinessTokenClaims{
		UserID:         claims.UserID,
		BusinessID:     claims.BusinessID,
		BusinessTypeID: claims.BusinessTypeID,
		RoleID:         claims.RoleID,
	}

	return &domain.AuthInfo{
		Type:                domain.AuthTypeJWT,
		UserID:              claims.UserID,
		BusinessID:          claims.BusinessID,
		BusinessTypeID:      claims.BusinessTypeID,
		RoleID:              claims.RoleID,
		JWTClaims:           claims,
		BusinessTokenClaims: businessTokenClaims,
	}, nil
}

// ValidateBusinessToken is now an alias for ValidateToken (for backward compatibility)
func (s *AuthService) ValidateBusinessToken(token string) (*domain.AuthInfo, error) {
	return s.ValidateToken(token)
}

// ValidateMainToken is now an alias for ValidateToken (for backward compatibility)
func (s *AuthService) ValidateMainToken(token string) (*domain.AuthInfo, error) {
	return s.ValidateToken(token)
}
