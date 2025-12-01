package domain

type JWTClaims struct {
	UserID    uint
	TokenType string
}
type BusinessTokenClaims struct {
	UserID         uint
	BusinessID     uint
	BusinessTypeID uint
	RoleID         uint
	TokenType      string
}
type ValidateAPIKeyRequest struct {
	APIKey string
}
type ValidateAPIKeyResponse struct {
	Success    bool
	Message    string
	UserID     uint
	Email      string
	BusinessID uint
	Roles      []string
	APIKeyID   uint
}
