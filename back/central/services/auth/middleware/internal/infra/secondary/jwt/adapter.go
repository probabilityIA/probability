package jwt

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware/internal/domain"
	sharedjwt "github.com/secamc93/probability/back/central/shared/jwt"
)

type Adapter struct {
	impl sharedjwt.IJWTService
}

func NewAdapter(impl sharedjwt.IJWTService) *Adapter {
	return &Adapter{impl: impl}
}

func (a *Adapter) GenerateToken(userID uint) (string, error) {
	return a.impl.GenerateToken(userID)
}

func (a *Adapter) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	claims, err := a.impl.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &domain.JWTClaims{
		UserID:    claims.UserID,
		TokenType: claims.TokenType,
	}, nil
}

func (a *Adapter) RefreshToken(tokenString string) (string, error) {
	return a.impl.RefreshToken(tokenString)
}

func (a *Adapter) GenerateBusinessToken(userID, businessID, businessTypeID, roleID uint) (string, error) {
	return a.impl.GenerateBusinessToken(userID, businessID, businessTypeID, roleID)
}

func (a *Adapter) ValidateBusinessToken(tokenString string) (*domain.BusinessTokenClaims, error) {
	claims, err := a.impl.ValidateBusinessToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &domain.BusinessTokenClaims{
		UserID:         claims.UserID,
		BusinessID:     claims.BusinessID,
		BusinessTypeID: claims.BusinessTypeID,
		RoleID:         claims.RoleID,
		TokenType:      claims.TokenType,
	}, nil
}
