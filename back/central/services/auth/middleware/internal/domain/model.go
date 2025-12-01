package domain

type AuthType string

const (
	AuthTypeUnknown AuthType = "unknown"
	AuthTypeJWT     AuthType = "jwt"
	AuthTypeAPIKey  AuthType = "api_key"
)

type AuthInfo struct {
	Type                AuthType
	UserID              uint
	Email               string
	Roles               []string
	BusinessID          uint
	APIKey              string
	JWTClaims           *JWTClaims
	BusinessTokenClaims *BusinessTokenClaims
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
