package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int64  `json:"created_at"`
}

func (c *JumpsellerClient) RefreshToken(ctx context.Context, tokenURL, clientID, clientSecret, refreshToken string) (*domain.TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("jumpseller client: creating refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, friendlyConnError(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("jumpseller client: reading refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d: %s", domain.ErrTokenRefreshFailed, resp.StatusCode, string(raw))
	}

	var parsed oauthTokenResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("jumpseller client: parsing refresh response: %w", err)
	}
	if parsed.AccessToken == "" {
		return nil, fmt.Errorf("%w: respuesta sin access_token", domain.ErrTokenRefreshFailed)
	}

	return &domain.TokenResponse{
		AccessToken:  parsed.AccessToken,
		RefreshToken: parsed.RefreshToken,
		TokenType:    parsed.TokenType,
		ExpiresIn:    parsed.ExpiresIn,
		CreatedAt:    parsed.CreatedAt,
	}, nil
}
