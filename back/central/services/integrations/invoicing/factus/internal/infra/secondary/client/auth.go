package client

import (
	"context"
	"fmt"
)

// AuthResponse representa la respuesta del endpoint OAuth de Factus
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`    // "Bearer"
	ExpiresIn    int    `json:"expires_in"`    // 600 (10 min) para access, 3600 para refresh
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

// authenticate obtiene un access token usando OAuth2 Password Grant o Refresh Token
// baseURL es opcional: si está vacío usa la baseURL del cliente HTTP
// Estrategia:
//  1. Si hay access_token válido en cache → retornar
//  2. Si access_token expirado pero refresh_token válido → usar refresh
//  3. Si ambos expirados → login completo con credenciales
func (c *Client) authenticate(ctx context.Context, baseURL, clientID, clientSecret, username, password string) (string, error) {
	// 1. Intentar usar access_token cacheado
	if token, ok := c.tokenCache.GetAccessToken(); ok {
		c.log.Debug(ctx).Msg("Using cached Factus access token")
		return token, nil
	}

	// 2. Intentar usar refresh_token
	if refreshToken, ok := c.tokenCache.GetRefreshToken(); ok {
		c.log.Info(ctx).Msg("🔄 Access token expired, attempting refresh...")
		token, err := c.refreshAccessToken(ctx, baseURL, clientID, clientSecret, refreshToken)
		if err == nil {
			return token, nil
		}
		c.log.Warn(ctx).Err(err).Msg("Refresh token failed, falling back to full login")
	}

	// 3. Login completo
	return c.loginWithPassword(ctx, baseURL, clientID, clientSecret, username, password)
}

// loginWithPassword realiza el login con credenciales completas
// POST /oauth/token con grant_type=password (form-data)
// baseURL es opcional: si está vacío usa la baseURL del cliente HTTP
func (c *Client) loginWithPassword(ctx context.Context, baseURL, clientID, clientSecret, username, password string) (string, error) {
	c.log.Info(ctx).
		Str("client_id_length", fmt.Sprintf("%d chars", len(clientID))).
		Str("username", username).
		Msg("🔑 Authenticating with Factus (password grant)")

	var authResp AuthResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"grant_type":    "password",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"username":      username,
			"password":      password,
		}).
		SetResult(&authResp).
		Post(c.endpointURL(baseURL, "/oauth/token"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Factus auth request failed - network error")
		return "", fmt.Errorf("error de red al conectar con Factus: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("token_type", authResp.TokenType).
		Int("expires_in", authResp.ExpiresIn).
		Msg("📥 Factus auth response received")

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Msg("❌ Factus authentication failed")
		switch resp.StatusCode() {
		case 401:
			return "", fmt.Errorf("credenciales inválidas: Client ID, Client Secret, usuario o contraseña incorrectos")
		case 422:
			return "", fmt.Errorf("datos de autenticación incompletos o con formato incorrecto")
		case 429:
			return "", fmt.Errorf("demasiadas solicitudes a Factus, intenta de nuevo en unos minutos")
		default:
			return "", fmt.Errorf("error de autenticación en Factus (código %d)", resp.StatusCode())
		}
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("Factus no retornó un token de acceso válido")
	}

	// Calcular TTL del refresh token (si viene expires_in es para refresh, access = 600s)
	// En práctica: access=600, refresh=3600
	accessTTL := 600
	refreshTTL := 3600
	if authResp.ExpiresIn > 0 {
		// El endpoint de password grant retorna expires_in del refresh_token
		refreshTTL = authResp.ExpiresIn
	}

	c.tokenCache.SetTokens(authResp.AccessToken, accessTTL, authResp.RefreshToken, refreshTTL)

	c.log.Info(ctx).
		Str("token_prefix", authResp.AccessToken[:min(20, len(authResp.AccessToken))]).
		Int("access_ttl_sec", accessTTL).
		Int("refresh_ttl_sec", refreshTTL).
		Msg("✅ Factus authentication successful")

	return authResp.AccessToken, nil
}

// refreshAccessToken obtiene un nuevo access_token usando el refresh_token
// POST /oauth/token con grant_type=refresh_token (form-data)
// baseURL es opcional: si está vacío usa la baseURL del cliente HTTP
func (c *Client) refreshAccessToken(ctx context.Context, baseURL, clientID, clientSecret, refreshToken string) (string, error) {
	c.log.Info(ctx).Msg("🔄 Refreshing Factus access token")

	var authResp AuthResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"refresh_token": refreshToken,
		}).
		SetResult(&authResp).
		Post(c.endpointURL(baseURL, "/oauth/token"))

	if err != nil {
		return "", fmt.Errorf("factus refresh request failed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("factus refresh failed (status %d)", resp.StatusCode())
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("factus refresh returned empty access_token")
	}

	// Refresh exitoso: nuevo access 10 min, nuevo refresh 1h
	accessTTL := 600
	refreshTTL := 3600
	if authResp.ExpiresIn > 0 {
		refreshTTL = authResp.ExpiresIn
	}

	c.tokenCache.SetTokens(authResp.AccessToken, accessTTL, authResp.RefreshToken, refreshTTL)

	c.log.Info(ctx).
		Str("token_prefix", authResp.AccessToken[:min(20, len(authResp.AccessToken))]).
		Msg("✅ Factus token refreshed successfully")

	return authResp.AccessToken, nil
}

// TestAuthentication valida las credenciales haciendo una autenticación de prueba
// baseURL es opcional: si está vacío usa la baseURL del cliente HTTP
func (c *Client) TestAuthentication(ctx context.Context, baseURL, clientID, clientSecret, username, password string) error {
	c.log.Info(ctx).
		Str("client_id_length", fmt.Sprintf("%d chars", len(clientID))).
		Str("username", username).
		Msg("🧪 Testing Factus credentials")

	// Limpiar cache para forzar nueva autenticación
	c.tokenCache.Clear()
	c.log.Info(ctx).Msg("🗑️ Token cache cleared - forcing fresh authentication")

	token, err := c.loginWithPassword(ctx, baseURL, clientID, clientSecret, username, password)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Factus credential validation failed")
		return err
	}

	if token == "" {
		return fmt.Errorf("Factus no retornó un token válido")
	}

	c.log.Info(ctx).
		Str("token_prefix", token[:min(20, len(token))]).
		Msg("✅ Factus credentials validated successfully")
	return nil
}

// min retorna el menor de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
