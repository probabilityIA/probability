package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// AuthRequest representa el request de autenticación a Softpymes
type AuthRequest struct {
	APIKey    string `json:"apiKey"`    // Softpymes API espera camelCase
	APISecret string `json:"apiSecret"` // Softpymes API espera camelCase
}

// AuthResponse representa la respuesta de autenticación de Softpymes
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`  // Token de acceso retornado por la API
	ExpiresInMin int    `json:"expiresInMin"` // Tiempo de expiración en minutos
	TokenType    string `json:"tokenType"`    // Tipo de token (Bearer)
	Success      bool   `json:"success"`      // Campo legacy, puede no venir
	Message      string `json:"message"`      // Mensaje de error si falla
	Error        string `json:"error"`        // Error detallado si falla
}

// authenticate obtiene un token de autenticación de Softpymes.
// baseURL: URL base efectiva (producción o testing); vacío usa c.baseURL.
// referer: Identificación de la instancia del cliente (requerido por API)
func (c *Client) authenticate(ctx context.Context, apiKey, apiSecret, referer, baseURL string) (string, error) {
	effectiveBase := baseURL
	if effectiveBase == "" {
		effectiveBase = c.baseURL
	}

	// Verificar si tenemos un token válido en cache para esta URL
	if token, valid := c.tokenCache.Get(effectiveBase); valid {
		c.log.Debug(ctx).Msg("Using cached authentication token")
		return token, nil
	}

	c.log.Info(ctx).
		Str("api_key_length", fmt.Sprintf("%d chars", len(apiKey))).
		Str("api_secret_length", fmt.Sprintf("%d chars", len(apiSecret))).
		Msg("🔑 Authenticating with Softpymes API")

	// Preparar request
	authReq := &AuthRequest{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	var authResp AuthResponse

	authURL := c.resolveURL(baseURL, "/oauth/integration/login/")

	c.log.Info(ctx).
		Str("endpoint", authURL).
		Interface("request_body", authReq).
		Msg("📤 Sending authentication request to Softpymes")

	// Hacer llamado a la API
	// Header Referer es requerido según documentación (identificación de la instancia del cliente)
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Referer", referer).
		SetBody(authReq).
		SetResult(&authResp).
		SetDebug(true).
		Post(authURL)

	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Msg("❌ Failed to authenticate with Softpymes - Network error")
		return "", fmt.Errorf("authentication request failed: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("status", resp.Status()).
		Interface("response_body", authResp).
		Msg("📥 Received authentication response from Softpymes")

	// Verificar respuesta
	if resp.IsError() {
		errorMsg := authResp.Error
		if errorMsg == "" {
			errorMsg = authResp.Message
		}
		if errorMsg == "" {
			// Intentar parsear el body raw para extraer mensaje de error
			// Maneja formatos como: {"error": {"code": "404", "message": "..."}}
			var genericError map[string]interface{}
			if err := json.Unmarshal(resp.Body(), &genericError); err == nil {
				// Intentar extraer error.message
				if errObj, ok := genericError["error"].(map[string]interface{}); ok {
					if msg, ok := errObj["message"].(string); ok && msg != "" {
						errorMsg = msg
					}
				}
			}
		}
		if errorMsg == "" {
			// Fallback al status HTTP
			errorMsg = fmt.Sprintf("HTTP %d - %s", resp.StatusCode(), resp.Status())
		}

		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", authResp.Error).
			Str("message", authResp.Message).
			Bool("success", authResp.Success).
			Str("final_error", errorMsg).
			Msg("❌ Authentication failed - HTTP error")
		return "", fmt.Errorf("autenticación falló (status %d): %s", resp.StatusCode(), errorMsg)
	}

	if authResp.AccessToken == "" {
		c.log.Error(ctx).
			Bool("success", authResp.Success).
			Str("message", authResp.Message).
			Str("error", authResp.Error).
			Int("token_length", len(authResp.AccessToken)).
			Msg("❌ Authentication unsuccessful - Empty access token")
		return "", fmt.Errorf("authentication failed: %s", authResp.Message)
	}

	// Guardar token en cache keyed por URL efectiva
	// Convertir minutos a segundos para el cache
	expiresInSeconds := authResp.ExpiresInMin * 60
	if expiresInSeconds == 0 {
		expiresInSeconds = 3600 // Default 1 hora si no viene el tiempo
	}
	c.tokenCache.Set(effectiveBase, authResp.AccessToken, expiresInSeconds)

	c.log.Info(ctx).
		Str("token_length", fmt.Sprintf("%d chars", len(authResp.AccessToken))).
		Int("expires_in_minutes", authResp.ExpiresInMin).
		Str("token_type", authResp.TokenType).
		Msg("✅ Successfully authenticated with Softpymes")

	return authResp.AccessToken, nil
}

// TestAuthentication valida las credenciales haciendo una autenticación de prueba.
// baseURL: URL base efectiva (producción o testing); vacío usa c.baseURL.
// referer: Identificación de la instancia del cliente (requerido por API)
func (c *Client) TestAuthentication(ctx context.Context, apiKey, apiSecret, referer, baseURL string) error {
	c.log.Info(ctx).
		Str("api_key_prefix", apiKey[:min(10, len(apiKey))]).
		Str("api_secret_length", fmt.Sprintf("%d chars", len(apiSecret))).
		Str("referer", referer).
		Msg("🧪 Testing Softpymes credentials")

	// Limpiar cache para forzar nueva autenticación
	c.tokenCache.Clear()
	c.log.Info(ctx).Msg("🗑️ Token cache cleared - forcing fresh authentication")

	// Intentar autenticar con todas las credenciales
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Msg("❌ Credential validation failed")
		return fmt.Errorf("credential validation failed: %w", err)
	}

	if token == "" {
		c.log.Error(ctx).Msg("❌ Authentication returned empty token")
		return fmt.Errorf("authentication returned empty token")
	}

	c.log.Info(ctx).
		Str("token_prefix", token[:min(20, len(token))]).
		Msg("✅ Credentials validated successfully")
	return nil
}

// Helper min function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
