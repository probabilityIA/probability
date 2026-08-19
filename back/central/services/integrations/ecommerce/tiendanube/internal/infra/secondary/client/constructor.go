package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

const (
	defaultUserAgent = "Probability (soporte@probabilityia.com.co)"
	maxRetries       = 3
	minInterval      = 500 * time.Millisecond
)

type TiendanubeClient struct {
	httpClient *http.Client

	mu       sync.Mutex
	lastCall map[string]time.Time
}

func New() domain.ITiendanubeClient {
	return &TiendanubeClient{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		lastCall:   make(map[string]time.Time),
	}
}

func (c *TiendanubeClient) pace(key string) {
	c.mu.Lock()
	last, ok := c.lastCall[key]
	now := time.Now()
	wait := time.Duration(0)
	if ok {
		if elapsed := now.Sub(last); elapsed < minInterval {
			wait = minInterval - elapsed
		}
	}
	c.lastCall[key] = now.Add(wait)
	c.mu.Unlock()

	if wait > 0 {
		time.Sleep(wait)
	}
}

func (c *TiendanubeClient) endpoint(cred domain.Credential, path string, query url.Values) (string, error) {
	base := strings.TrimRight(strings.TrimSpace(cred.BaseURL), "/")
	if base == "" {
		return "", domain.ErrMissingBaseURL
	}
	storeID := strings.TrimSpace(cred.StoreID)
	if storeID == "" {
		return "", domain.ErrMissingStoreID
	}
	if !strings.HasSuffix(base, "/"+storeID) {
		base = base + "/" + storeID
	}
	full := base + path
	if len(query) > 0 {
		full = full + "?" + query.Encode()
	}
	return full, nil
}

func userAgent(cred domain.Credential) string {
	if ua := strings.TrimSpace(cred.UserAgent); ua != "" {
		return ua
	}
	return defaultUserAgent
}

func (c *TiendanubeClient) do(ctx context.Context, cred domain.Credential, method, path string, query url.Values, body interface{}) ([]byte, http.Header, error) {
	endpoint, err := c.endpoint(cred, path, query)
	if err != nil {
		return nil, nil, err
	}

	var payload []byte
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("tiendanube client: encoding body: %w", err)
		}
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		c.pace(cred.AccessToken)

		var reader io.Reader
		if payload != nil {
			reader = bytes.NewReader(payload)
		}

		req, rerr := http.NewRequestWithContext(ctx, method, endpoint, reader)
		if rerr != nil {
			return nil, nil, fmt.Errorf("tiendanube client: creating request: %w", rerr)
		}
		req.Header.Set("Authentication", "bearer "+cred.AccessToken)
		req.Header.Set("Authorization", "Bearer "+cred.AccessToken)
		req.Header.Set("User-Agent", userAgent(cred))
		req.Header.Set("Accept", "application/json")
		if payload != nil {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}

		resp, derr := c.httpClient.Do(req)
		if derr != nil {
			return nil, nil, fmt.Errorf("tiendanube client: request failed: %w", derr)
		}

		raw, _ := io.ReadAll(resp.Body)
		headers := resp.Header
		status := resp.StatusCode
		resp.Body.Close()

		switch {
		case status == http.StatusUnauthorized || status == http.StatusForbidden:
			return nil, headers, domain.ErrInvalidCredentials
		case status == http.StatusTooManyRequests:
			if attempt == maxRetries-1 {
				return nil, headers, domain.ErrRateLimited
			}
			time.Sleep(retryAfter(headers, attempt))
			continue
		case status >= 500:
			if attempt == maxRetries-1 {
				return nil, headers, fmt.Errorf("tiendanube client: %s %s returned %d: %s", method, path, status, truncate(raw))
			}
			time.Sleep(retryAfter(headers, attempt))
			continue
		case status >= 400:
			return nil, headers, fmt.Errorf("tiendanube client: %s %s returned %d: %s", method, path, status, truncate(raw))
		}

		return raw, headers, nil
	}

	return nil, nil, fmt.Errorf("tiendanube client: %s %s exhausted retries", method, path)
}

func retryAfter(headers http.Header, attempt int) time.Duration {
	if headers != nil {
		if ms := headers.Get("x-rate-limit-reset"); ms != "" {
			if n, err := strconv.Atoi(ms); err == nil && n > 0 && n < 60000 {
				return time.Duration(n) * time.Millisecond
			}
		}
	}
	return time.Duration(attempt+1) * time.Second
}

func truncate(raw []byte) string {
	const limit = 300
	if len(raw) <= limit {
		return string(raw)
	}
	return string(raw[:limit]) + "..."
}
