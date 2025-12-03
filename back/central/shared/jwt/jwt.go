package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IJWTService define operaciones de JWT sin depender de otros módulos
type IJWTService interface {
	// Token unificado que incluye toda la información
	GenerateToken(userID, businessID, businessTypeID, roleID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(tokenString string) (string, error)

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

// Claims representa los claims internos del token unificado
type Claims struct {
	UserID         uint `json:"user_id"`
	BusinessID     uint `json:"business_id"`
	BusinessTypeID uint `json:"business_type_id"`
	RoleID         uint `json:"role_id"`
	jwt.RegisteredClaims
}

// JWTClaims es la estructura pública que exponemos a consumidores
type JWTClaims struct {
	UserID         uint
	BusinessID     uint
	BusinessTypeID uint
	RoleID         uint
}

// New crea una nueva instancia del servicio JWT (autocontenida)
func New(secretKey string) IJWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}

// GenerateToken genera un nuevo token JWT unificado con toda la información
func (j *JWTService) GenerateToken(userID, businessID, businessTypeID, roleID uint) (string, error) {
	claims := Claims{
		UserID:         userID,
		BusinessID:     businessID,
		BusinessTypeID: businessTypeID,
		RoleID:         roleID,
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
			UserID:         claims.UserID,
			BusinessID:     claims.BusinessID,
			BusinessTypeID: claims.BusinessTypeID,
			RoleID:         claims.RoleID,
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

	return j.GenerateToken(claims.UserID, claims.BusinessID, claims.BusinessTypeID, claims.RoleID)
}
