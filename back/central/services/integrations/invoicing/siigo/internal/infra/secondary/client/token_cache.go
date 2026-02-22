package client

import (
	"sync"
	"time"
)

// TokenCache maneja el cache del token de autenticación de Siigo
// Siigo tiene TTL de access_token = 24h (86400 segundos), sin refresh token
type TokenCache struct {
	accessToken        string
	accessTokenExpires time.Time
	mu                 sync.RWMutex
}

// NewTokenCache crea un nuevo cache de token para Siigo
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

// SetToken guarda el access_token en el cache
// expiresIn: TTL del token en segundos (86400 = 24h)
func (tc *TokenCache) SetToken(accessToken string, expiresIn int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Buffer de 30 min para el token (efectivo 23.5h de 24h)
	buffer := 1800
	ttl := expiresIn - buffer
	if ttl <= 0 {
		ttl = expiresIn / 2
	}

	tc.accessToken = accessToken
	tc.accessTokenExpires = time.Now().Add(time.Duration(ttl) * time.Second)
}

// Clear limpia el token del cache
func (tc *TokenCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.accessToken = ""
	tc.accessTokenExpires = time.Time{}
}
