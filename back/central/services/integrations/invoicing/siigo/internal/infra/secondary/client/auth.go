package client

import (
	"context"
	"fmt"
)

// AuthRequest representa el body del endpoint de autenticación de Siigo
type AuthRequest struct {
	Username  string `json:"username"`
	AccessKey string `json:"access_key"`
}

// AuthResponse representa la respuesta del endpoint de autenticación de Siigo
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`    // 86400 (24h)
	TokenType   string `json:"token_type"`    // "Bearer"
}

// authenticate obtiene un access token de Siigo
// Estrategia:
//  1. Si hay access_token válido en cache → retornar
//  2. Si expirado → login completo con credenciales (no hay refresh en Siigo)
func (c *Client) authenticate(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) (string, error) {
	// 1. Intentar usar access_token cacheado
	if token, ok := c.tokenCache.GetAccessToken(); ok {
		c.log.Debug(ctx).Msg("Using cached Siigo access token")
		return token, nil
	}

	// 2. Login completo (Siigo no tiene refresh token)
	return c.loginWithCredentials(ctx, username, accessKey, accountID, partnerID, baseURL)
}

// loginWithCredentials realiza el login con credenciales de Siigo
// POST /v1/auth con JSON body y headers especiales
func (c *Client) loginWithCredentials(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) (string, error) {
	c.log.Info(ctx).
		Str("username", username).
		Msg("🔑 Authenticating with Siigo")

	var authResp AuthResponse

	authBody := AuthRequest{
		Username:  username,
		AccessKey: accessKey,
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", accountID).
		SetHeader("Partner-Id", partnerID).
		SetBody(authBody).
		SetResult(&authResp).
		Post(c.endpointURL(baseURL, "/v1/auth"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo auth request failed - network error")
		return "", fmt.Errorf("error de red al conectar con Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("token_type", authResp.TokenType).
		Int("expires_in", authResp.ExpiresIn).
		Msg("📥 Siigo auth response received")

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Siigo authentication failed")
		switch resp.StatusCode() {
		case 401:
			return "", fmt.Errorf("credenciales inválidas: username o access_key incorrectos")
		case 403:
			return "", fmt.Errorf("acceso denegado: account_id o partner_id inválidos")
		case 422:
			return "", fmt.Errorf("datos de autenticación incompletos o con formato incorrecto")
		case 429:
			return "", fmt.Errorf("demasiadas solicitudes a Siigo, intenta de nuevo en unos minutos")
		default:
			return "", fmt.Errorf("error de autenticación en Siigo (código %d)", resp.StatusCode())
		}
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("Siigo no retornó un token de acceso válido")
	}

	// TTL por defecto 86400 si la respuesta no trae expires_in
	ttl := authResp.ExpiresIn
	if ttl <= 0 {
		ttl = 86400
	}

	c.tokenCache.SetToken(authResp.AccessToken, ttl)

	tokenPreview := authResp.AccessToken
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}

	c.log.Info(ctx).
		Str("token_prefix", tokenPreview).
		Int("expires_in_sec", ttl).
		Msg("✅ Siigo authentication successful")

	return authResp.AccessToken, nil
}

// TestAuthentication valida las credenciales haciendo una autenticación de prueba
func (c *Client) TestAuthentication(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error {
	c.log.Info(ctx).
		Str("username", username).
		Msg("🧪 Testing Siigo credentials")

	// Limpiar cache para forzar nueva autenticación
	c.tokenCache.Clear()
	c.log.Info(ctx).Msg("🗑️ Token cache cleared - forcing fresh authentication")

	token, err := c.loginWithCredentials(ctx, username, accessKey, accountID, partnerID, baseURL)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo credential validation failed")
		return err
	}

	if token == "" {
		return fmt.Errorf("Siigo no retornó un token válido")
	}

	tokenPreview := token
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}

	c.log.Info(ctx).
		Str("token_prefix", tokenPreview).
		Msg("✅ Siigo credentials validated successfully")
	return nil
}
