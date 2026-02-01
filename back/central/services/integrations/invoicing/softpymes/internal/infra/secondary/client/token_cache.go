package client

import (
	"sync"
	"time"
)

// TokenCache maneja el cache de tokens de autenticación
type TokenCache struct {
	token     string
	expiresAt time.Time
	mu        sync.RWMutex
}

// NewTokenCache crea un nuevo cache de tokens
func NewTokenCache() *TokenCache {
	return &TokenCache{}
}

// Get obtiene el token del cache si es válido
func (tc *TokenCache) Get() (string, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if tc.token == "" || time.Now().After(tc.expiresAt) {
		return "", false
	}

	return tc.token, true
}

// Set guarda un token en el cache
func (tc *TokenCache) Set(token string, expiresIn int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Restar 5 minutos al tiempo de expiración para renovar antes
	expirationTime := time.Now().Add(time.Duration(expiresIn-300) * time.Second)

	tc.token = token
	tc.expiresAt = expirationTime
}

// Clear limpia el cache
func (tc *TokenCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.token = ""
	tc.expiresAt = time.Time{}
}
