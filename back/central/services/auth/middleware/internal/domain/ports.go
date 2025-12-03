package domain

import "context"

type IJWTService interface {
	// Token unificado que incluye toda la información
	GenerateToken(userID, businessID, businessTypeID, roleID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)
}

type IAuthUseCase interface {
	ValidateAPIKey(ctx context.Context, request ValidateAPIKeyRequest) (*ValidateAPIKeyResponse, error)
}
