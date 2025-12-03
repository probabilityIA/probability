package domain

type JWTClaims struct {
	UserID         uint
	BusinessID     uint
	BusinessTypeID uint
	RoleID         uint
}
type BusinessTokenClaims struct {
	UserID         uint
	BusinessID     uint
	BusinessTypeID uint
	RoleID         uint
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
