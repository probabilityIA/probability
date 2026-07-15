package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

const meliAPIBaseURL = "https://api.mercadolibre.com"

const (
	maxRetries429  = 5
	initialBackoff = 1 * time.Second
	maxBackoff     = 60 * time.Second
)

type MeliClient struct {
	httpClient *http.Client
	baseURL    string
	limiter    *rateLimiter
}

func New() domain.IMeliClient {
	return &MeliClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: meliAPIBaseURL,
		limiter: newRateLimiter(100),
	}
}

func (c *MeliClient) newAuthorizedRequest(ctx context.Context, method, url, accessToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("meli client: creating request: %w", err)
	}
	c.setHeaders(req, accessToken)
	return req, nil
}

func (c *MeliClient) newAuthorizedRequestWithBody(ctx context.Context, method, url, accessToken string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("meli client: creating request: %w", err)
	}
	c.setHeaders(req, accessToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *MeliClient) setHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("x-format-new", "true")
}

func (c *MeliClient) do(ctx context.Context, build func() (*http.Request, error)) (*http.Response, []byte, error) {
	backoff := initialBackoff
	var lastErr error
	for attempt := 0; attempt <= maxRetries429; attempt++ {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, nil, err
		}
		req, err := build()
		if err != nil {
			return nil, nil, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("meli client: request failed: %w", err)
			if !sleepBackoff(ctx, &backoff) {
				return nil, nil, ctx.Err()
			}
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = domain.ErrRateLimited
			if attempt == maxRetries429 {
				break
			}
			if !sleepBackoff(ctx, &backoff) {
				return nil, nil, ctx.Err()
			}
			continue
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return resp, nil, fmt.Errorf("meli client: reading response: %w", readErr)
		}
		return resp, body, nil
	}
	return nil, nil, lastErr
}

func sleepBackoff(ctx context.Context, backoff *time.Duration) bool {
	timer := time.NewTimer(*backoff)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
	}
	*backoff *= 2
	if *backoff > maxBackoff {
		*backoff = maxBackoff
	}
	return true
}

func (c *MeliClient) TestConnection(ctx context.Context, accessToken string) error {
	resp, _, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, fmt.Sprintf("%s/users/me", c.baseURL), accessToken)
	})
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("meli client: unexpected status %d", resp.StatusCode)
	}
	return nil
}
