package client

import (
	"sync"
	"time"
)

// TokenCache maneja el cache de tokens de autenticación de Factus
// Factus tiene TTL de access_token = 10 min, refresh_token = 1h
type TokenCache struct {
	accessToken        string
	accessTokenExpires time.Time
	refreshToken       string
	refreshTokenExpires time.Time
	mu                 sync.RWMutex
}

// NewTokenCache crea un nuevo cache de tokens para Factus
func NewTokenCache() *TokenCache {
	return &TokenCache{}
}

// GetAccessToken obtiene el access token del cache si es válido
func (tc *TokenCache) GetAccessToken() (string, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if tc.accessToken == "" || time.Now().After(tc.accessTokenExpires) {
		return "", false
	}
	return tc.accessToken, true
}

// GetRefreshToken obtiene el refresh token del cache si es válido
func (tc *TokenCache) GetRefreshToken() (string, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if tc.refreshToken == "" || time.Now().After(tc.refreshTokenExpires) {
		return "", false
	}
	return tc.refreshToken, true
}

// SetTokens guarda access_token y refresh_token en el cache
// accessExpiresIn: TTL del access token en segundos (600 = 10 min)
// refreshExpiresIn: TTL del refresh token en segundos (3600 = 1h)
func (tc *TokenCache) SetTokens(accessToken string, accessExpiresIn int, refreshToken string, refreshExpiresIn int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Buffer de 2 min para access token (efectivo 8 min de 10 min)
	accessBuffer := 120
	accessTTL := accessExpiresIn - accessBuffer
	if accessTTL <= 0 {
		accessTTL = accessExpiresIn / 2
	}

	// Buffer de 5 min para refresh token (efectivo 55 min de 60 min)
	refreshBuffer := 300
	refreshTTL := refreshExpiresIn - refreshBuffer
	if refreshTTL <= 0 {
		refreshTTL = refreshExpiresIn / 2
	}

	tc.accessToken = accessToken
	tc.accessTokenExpires = time.Now().Add(time.Duration(accessTTL) * time.Second)

	if refreshToken != "" {
		tc.refreshToken = refreshToken
		tc.refreshTokenExpires = time.Now().Add(time.Duration(refreshTTL) * time.Second)
	}
}

// Clear limpia todos los tokens del cache
func (tc *TokenCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.accessToken = ""
	tc.accessTokenExpires = time.Time{}
	tc.refreshToken = ""
	tc.refreshTokenExpires = time.Time{}
}
