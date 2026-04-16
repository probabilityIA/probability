package client

import (
	"sync"
	"time"
)

// tokenEntry almacena un token y su tiempo de expiración
type tokenEntry struct {
	token     string
	expiresAt time.Time
}

// TokenCache maneja el cache de tokens de autenticación keyed por baseURL.
// Permite usar tokens distintos para producción y testing simultáneamente.
type TokenCache struct {
	mu      sync.Mutex
	entries map[string]*tokenEntry
}

// NewTokenCache crea un nuevo cache de tokens
func NewTokenCache() *TokenCache {
	return &TokenCache{
		entries: make(map[string]*tokenEntry),
	}
}

// Get obtiene el token del cache si es válido para la baseURL dada
func (tc *TokenCache) Get(baseURL string) (string, bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	entry, ok := tc.entries[baseURL]
	if !ok || entry.token == "" || time.Now().After(entry.expiresAt) {
		return "", false
	}

	return entry.token, true
}

// Set guarda un token en el cache para la baseURL dada
func (tc *TokenCache) Set(baseURL, token string, expiresIn int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Restar 5 minutos al tiempo de expiración para renovar antes
	expirationTime := time.Now().Add(time.Duration(expiresIn-300) * time.Second)

	tc.entries[baseURL] = &tokenEntry{
		token:     token,
		expiresAt: expirationTime,
	}
}

// Clear limpia todos los tokens del cache
func (tc *TokenCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.entries = make(map[string]*tokenEntry)
}
