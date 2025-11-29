package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IJWTService define operaciones de JWT sin depender de otros módulos
type IJWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)

	// Tokens para business
	GenerateBusinessToken(userID, businessID, businessTypeID, roleID uint) (string, error)
	ValidateBusinessToken(tokenString string) (*BusinessTokenClaims, error)

	// Tokens para votación pública
	GeneratePublicVotingToken(votingID, votingGroupID, hpID uint, durationHours int) (string, error)
	GenerateVotingAuthToken(residentID, propertyUnitID, votingID, votingGroupID, hpID uint) (string, error)
	ValidatePublicVotingToken(tokenString string) (*PublicVotingClaims, error)
	ValidateVotingAuthToken(tokenString string) (*VotingAuthClaims, error)
}

// JWTService implementación concreta
type JWTService struct {
	secretKey string
}

// Claims representa los claims internos del token
type Claims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"` // "main" o "business"
	jwt.RegisteredClaims
}

// JWTClaims es la estructura pública que exponemos a consumidores
type JWTClaims struct {
	UserID    uint
	TokenType string
}

// BusinessTokenClaims representa los claims del token de business
type BusinessTokenClaims struct {
	UserID         uint
	BusinessID     uint
	BusinessTypeID uint
	RoleID         uint
	TokenType      string
}

// BusinessClaims representa los claims internos del business token
type BusinessClaims struct {
	UserID         uint   `json:"user_id"`
	BusinessID     uint   `json:"business_id"`
	BusinessTypeID uint   `json:"business_type_id"`
	RoleID         uint   `json:"role_id"`
	TokenType      string `json:"token_type"`
	jwt.RegisteredClaims
}

// New crea una nueva instancia del servicio JWT (autocontenida)
func New(secretKey string) IJWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}

// GenerateToken genera un nuevo token JWT
func (j *JWTService) GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: "main",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "central-reserve-api",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("error al firmar token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida y decodifica un token JWT
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &JWTClaims{
			UserID:    claims.UserID,
			TokenType: claims.TokenType,
		}, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// RefreshToken refresca un token JWT
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return j.GenerateToken(claims.UserID)
}

// GenerateBusinessToken genera un token JWT para un business específico
func (j *JWTService) GenerateBusinessToken(userID, businessID, businessTypeID, roleID uint) (string, error) {
	claims := BusinessClaims{
		UserID:         userID,
		BusinessID:     businessID,
		BusinessTypeID: businessTypeID,
		RoleID:         roleID,
		TokenType:      "business",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "central-reserve-api",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("error al firmar business token: %w", err)
	}

	return tokenString, nil
}

// ValidateBusinessToken valida y decodifica un business token JWT
func (j *JWTService) ValidateBusinessToken(tokenString string) (*BusinessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &BusinessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear business token: %w", err)
	}

	if claims, ok := token.Claims.(*BusinessClaims); ok && token.Valid {
		return &BusinessTokenClaims{
			UserID:         claims.UserID,
			BusinessID:     claims.BusinessID,
			BusinessTypeID: claims.BusinessTypeID,
			RoleID:         claims.RoleID,
			TokenType:      claims.TokenType,
		}, nil
	}

	return nil, fmt.Errorf("business token inválido")
}
