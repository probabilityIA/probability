package domain

import "context"

type IJWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)

	// Tokens para business
	GenerateBusinessToken(userID, businessID, businessTypeID, roleID uint) (string, error)
	ValidateBusinessToken(tokenString string) (*BusinessTokenClaims, error)
}

type IAuthUseCase interface {
	ValidateAPIKey(ctx context.Context, request ValidateAPIKeyRequest) (*ValidateAPIKeyResponse, error)
}
