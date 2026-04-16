package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PublicVotingClaims - Claims para tokens de votación pública
type PublicVotingClaims struct {
	VotingID      uint   `json:"voting_id"`
	VotingGroupID uint   `json:"voting_group_id"`
	HPID          uint   `json:"hp_id"`
	Scope         string `json:"scope"` // "public_voting"
	jwt.RegisteredClaims
}

// VotingAuthClaims - Claims después de validar residente
type VotingAuthClaims struct {
	ResidentID     uint   `json:"resident_id"`
	PropertyUnitID uint   `json:"property_unit_id"`
	VotingID       uint   `json:"voting_id"`
	VotingGroupID  uint   `json:"voting_group_id"`
	HPID           uint   `json:"hp_id"`
	Scope          string `json:"scope"` // "voting_auth"
	jwt.RegisteredClaims
}

// GeneratePublicVotingToken genera un token para acceder a la página pública de votación
func (j *JWTService) GeneratePublicVotingToken(votingID, votingGroupID, hpID uint, durationHours int) (string, error) {
	if durationHours <= 0 {
		durationHours = 24 // Default 24 horas
	}

	claims := PublicVotingClaims{
		VotingID:      votingID,
		VotingGroupID: votingGroupID,
		HPID:          hpID,
		Scope:         "public_voting",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(durationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("public_voting_%d", votingID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("error generando token de votación pública: %w", err)
	}

	return tokenString, nil
}

// GenerateVotingAuthToken genera un token temporal después de validar al residente
func (j *JWTService) GenerateVotingAuthToken(residentID, propertyUnitID, votingID, votingGroupID, hpID uint) (string, error) {
	claims := VotingAuthClaims{
		ResidentID:     residentID,
		PropertyUnitID: propertyUnitID,
		VotingID:       votingID,
		VotingGroupID:  votingGroupID,
		HPID:           hpID,
		Scope:          "voting_auth",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // 2 horas para votar
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("voting_auth_%d_%d", residentID, votingID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("error generando token de autenticación de votación: %w", err)
	}

	return tokenString, nil
}

// ValidatePublicVotingToken valida un token de votación pública
func (j *JWTService) ValidatePublicVotingToken(tokenString string) (*PublicVotingClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PublicVotingClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if claims, ok := token.Claims.(*PublicVotingClaims); ok && token.Valid {
		if claims.Scope != "public_voting" {
			return nil, fmt.Errorf("scope inválido para token de votación pública")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("token de votación pública inválido")
}

// ValidateVotingAuthToken valida un token de autenticación de votación
func (j *JWTService) ValidateVotingAuthToken(tokenString string) (*VotingAuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &VotingAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if claims, ok := token.Claims.(*VotingAuthClaims); ok && token.Valid {
		if claims.Scope != "voting_auth" {
			return nil, fmt.Errorf("scope inválido para token de autenticación de votación")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("token de autenticación de votación inválido")
}
