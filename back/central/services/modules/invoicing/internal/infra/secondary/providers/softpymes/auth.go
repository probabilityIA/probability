package softpymes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/response"
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

// Authenticate obtiene un token de autenticación de Softpymes
func (c *Client) Authenticate(ctx context.Context, credentials map[string]interface{}) (string, error) {
	// Verificar si tenemos un token válido en cache
	if token, valid := c.tokenCache.Get(); valid {
		c.log.Info(ctx).Msg("Using cached authentication token")
		return token, nil
	}

	c.log.Info(ctx).Msg("Authenticating with Softpymes API")

	// Extraer credenciales
	apiKey, ok := credentials["api_key"].(string)
	if !ok || apiKey == "" {
		return "", fmt.Errorf("missing or invalid api_key")
	}

	apiSecret, ok := credentials["api_secret"].(string)
	if !ok || apiSecret == "" {
		return "", fmt.Errorf("missing or invalid api_secret")
	}

	// Preparar request
	authReq := &request.AuthRequest{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	var authResp response.AuthResponse

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(authReq).
		SetResult(&authResp).
		Post("/get_token")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to authenticate with Softpymes")
		return "", fmt.Errorf("authentication request failed: %w", err)
	}

	// Verificar respuesta
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", authResp.Error).
			Msg("Authentication failed")
		return "", fmt.Errorf("authentication failed: %s", authResp.Error)
	}

	if !authResp.Success || authResp.Token == "" {
		return "", fmt.Errorf("authentication failed: %s", authResp.Message)
	}

	// Guardar token en cache
	c.tokenCache.Set(authResp.Token, authResp.ExpiresIn)

	c.log.Info(ctx).
		Int("expires_in", authResp.ExpiresIn).
		Msg("Successfully authenticated with Softpymes")

	return authResp.Token, nil
}

// ValidateCredentials valida las credenciales haciendo una autenticación de prueba
func (c *Client) ValidateCredentials(ctx context.Context, credentials map[string]interface{}) error {
	// Limpiar cache para forzar nueva autenticación
	c.tokenCache.Clear()

	// Intentar autenticar
	token, err := c.Authenticate(ctx, credentials)
	if err != nil {
		return err
	}

	if token == "" {
		return fmt.Errorf("authentication returned empty token")
	}

	c.log.Info(ctx).Msg("Credentials validated successfully")
	return nil
}
