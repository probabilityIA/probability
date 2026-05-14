package client

import (
	"context"
	"fmt"
)

type AuthRequest struct {
	Username  string `json:"username"`
	AccessKey string `json:"access_key"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (c *Client) authenticate(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) (string, error) {
	if token, ok := c.tokenStore.Get(username, accountID, partnerID, baseURL); ok {
		c.log.Debug(ctx).Msg("Using cached Siigo access token")
		return token, nil
	}

	return c.loginWithCredentials(ctx, username, accessKey, accountID, partnerID, baseURL)
}

func (c *Client) loginWithCredentials(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) (string, error) {
	c.log.Info(ctx).
		Str("username", username).
		Msg("Authenticating with Siigo")

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
		c.log.Error(ctx).Err(err).Msg("Siigo auth request failed - network error")
		return "", fmt.Errorf("error de red al conectar con Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Str("token_type", authResp.TokenType).
		Int("expires_in", authResp.ExpiresIn).
		Msg("Siigo auth response received")

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo authentication failed")
		switch resp.StatusCode() {
		case 401:
			return "", fmt.Errorf("credenciales invalidas: username o access_key incorrectos")
		case 403:
			return "", fmt.Errorf("acceso denegado: account_id o partner_id invalidos")
		case 422:
			return "", fmt.Errorf("datos de autenticacion incompletos o con formato incorrecto")
		case 429:
			return "", fmt.Errorf("demasiadas solicitudes a Siigo, intenta de nuevo en unos minutos")
		default:
			return "", fmt.Errorf("error de autenticacion en Siigo (codigo %d)", resp.StatusCode())
		}
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("Siigo no retorno un token de acceso valido")
	}

	ttl := authResp.ExpiresIn
	if ttl <= 0 {
		ttl = 86400
	}

	c.tokenStore.Set(username, accountID, partnerID, baseURL, authResp.AccessToken, ttl)

	tokenPreview := authResp.AccessToken
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}

	c.log.Info(ctx).
		Str("token_prefix", tokenPreview).
		Int("expires_in_sec", ttl).
		Msg("Siigo authentication successful")

	return authResp.AccessToken, nil
}

func (c *Client) TestAuthentication(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error {
	c.log.Info(ctx).
		Str("username", username).
		Msg("Testing Siigo credentials")

	c.tokenStore.Clear(username, accountID, partnerID, baseURL)
	c.log.Info(ctx).Msg("Token cache cleared - forcing fresh authentication")

	token, err := c.loginWithCredentials(ctx, username, accessKey, accountID, partnerID, baseURL)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo credential validation failed")
		return err
	}

	if token == "" {
		return fmt.Errorf("Siigo no retorno un token valido")
	}

	tokenPreview := token
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}

	c.log.Info(ctx).
		Str("token_prefix", tokenPreview).
		Msg("Siigo credentials validated successfully")
	return nil
}
