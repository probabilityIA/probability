package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/client/response"
)

func (c *TikTokClient) GetClientToken(ctx context.Context, baseURL, clientKey, clientSecret string) (string, error) {
	endpoint := strings.TrimSuffix(baseURL, "/") + "/v2/oauth/token/"

	form := url.Values{}
	form.Set("client_key", clientKey)
	form.Set("client_secret", clientSecret)
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("tiktok: building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("tiktok: token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("tiktok: reading token response: %w", err)
	}

	var tokenResp response.TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("tiktok: parsing token response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK || tokenResp.AccessToken == "" {
		if tokenResp.Error != "" {
			return "", fmt.Errorf("tiktok: token error %s: %s", tokenResp.Error, tokenResp.ErrorDescription)
		}
		return "", fmt.Errorf("tiktok: token request returned status %d", resp.StatusCode)
	}

	return tokenResp.AccessToken, nil
}
