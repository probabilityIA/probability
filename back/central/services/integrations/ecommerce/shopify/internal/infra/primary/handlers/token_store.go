package handlers

import (
	"sync"
	"time"
)

// TokenExchangeData almacena las credenciales temporalmente para el intercambio
type TokenExchangeData struct {
	AccessToken  string
	ClientID     string
	ClientSecret string
	Expiry       time.Time
}

// TokenStore gestiona el almacenamiento temporal de tokens de intercambio
// Usamos un mapa simple con mutex para thread-safety
var (
	tokenStore = make(map[string]TokenExchangeData)
	tokenMutex sync.RWMutex
)

// StoreExchangeToken guarda un token temporal
func StoreExchangeToken(token string, data TokenExchangeData) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// Limpieza lazy: eliminar tokens expirados al insertar uno nuevo (simple)
	// En prod idealmente usaríamos un goroutine de limpieza o Redis
	now := time.Now()
	for k, v := range tokenStore {
		if now.After(v.Expiry) {
			delete(tokenStore, k)
		}
	}

	tokenStore[token] = data
}

// RetrieveExchangeToken recupera y elimina (consume) un token temporal
func RetrieveExchangeToken(token string) (TokenExchangeData, bool) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	data, exists := tokenStore[token]
	if !exists {
		return TokenExchangeData{}, false
	}

	// Verificar expiración
	if time.Now().After(data.Expiry) {
		delete(tokenStore, token)
		return TokenExchangeData{}, false
	}

	// Consumir el token (one-time use)
	delete(tokenStore, token)
	return data, true
}
