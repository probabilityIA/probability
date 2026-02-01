package client

import (
	"context"
	"fmt"
)

// AuthRequest representa el request de autenticación a Softpymes
type AuthRequest struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// AuthResponse representa la respuesta de autenticación de Softpymes
type AuthResponse struct {
	Success   bool   `json:"success"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}

// authenticate obtiene un token de autenticación de Softpymes
func (c *Client) authenticate(ctx context.Context, apiKey, apiSecret string) (string, error) {
	// Verificar si tenemos un token válido en cache
	if token, valid := c.tokenCache.Get(); valid {
		c.log.Debug(ctx).Msg("Using cached authentication token")
		return token, nil
	}

	c.log.Info(ctx).Msg("Authenticating with Softpymes API")

	// Preparar request
	authReq := &AuthRequest{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	var authResp AuthResponse

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

// TestAuthentication valida las credenciales haciendo una autenticación de prueba
func (c *Client) TestAuthentication(ctx context.Context, apiKey string) error {
	c.log.Info(ctx).Msg("Testing Softpymes credentials")

	// Limpiar cache para forzar nueva autenticación
	c.tokenCache.Clear()

	// Por ahora, apiKey contiene tanto la key como el secret separados por ":"
	// En producción, esto debería venir del map de credentials
	// Formato esperado: "api_key:api_secret"

	// Para simplificar, asumimos que tenemos solo api_key por ahora
	// En la implementación real, necesitaríamos api_secret también
	apiSecret := "" // TODO: Obtener de credentials map

	// Intentar autenticar
	token, err := c.authenticate(ctx, apiKey, apiSecret)
	if err != nil {
		return err
	}

	if token == "" {
		return fmt.Errorf("authentication returned empty token")
	}

	c.log.Info(ctx).Msg("Credentials validated successfully")
	return nil
}
