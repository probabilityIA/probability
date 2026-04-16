package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

// refreshTokenRequest es el body para POST /oauth/token.
type refreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken renueva el access_token usando el refresh_token.
// POST https://api.mercadolibre.com/oauth/token
func (c *MeliClient) RefreshToken(ctx context.Context, appID, clientSecret, refreshToken string) (*domain.TokenResponse, error) {
	endpoint := fmt.Sprintf("%s/oauth/token", c.baseURL)

	reqBody := refreshTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     appID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("meli client: marshaling refresh request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("meli client: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("meli client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", domain.ErrTokenRefreshFailed, resp.StatusCode, string(body))
	}

	var tokenResp response.MeliTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("meli client: parsing token response: %w", err)
	}

	token := tokenResp.ToDomain()
	return &token, nil
}
