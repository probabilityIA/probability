package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

type VTEXClient struct {
	httpClient *http.Client
}

func New() domain.IVTEXClient {
	return &VTEXClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func baseURL(cred domain.Credential) string {
	return fmt.Sprintf("https://%s.vtexcommercestable.com.br", cred.AccountName)
}

func (c *VTEXClient) newRequest(ctx context.Context, method, url string, cred domain.Credential, body []byte) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("vtex client: creating request: %w", err)
	}
	req.Header.Set("X-VTEX-API-AppKey", cred.AppKey)
	req.Header.Set("X-VTEX-API-AppToken", cred.AppToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *VTEXClient) do(ctx context.Context, method, url string, cred domain.Credential, body []byte) ([]byte, error) {
	req, err := c.newRequest(ctx, method, url, cred, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vtex client: request failed: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vtex client: reading response: %w", err)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, domain.ErrInvalidCredentials
	case resp.StatusCode == http.StatusTooManyRequests:
		return nil, domain.ErrRateLimited
	case resp.StatusCode == http.StatusNotFound:
		return nil, domain.ErrProductNotFound
	case resp.StatusCode >= 400:
		return nil, fmt.Errorf("vtex client: status %d: %s", resp.StatusCode, string(payload))
	}

	return payload, nil
}

func (c *VTEXClient) TestConnection(ctx context.Context, cred domain.Credential) error {
	endpoint := fmt.Sprintf("%s/api/catalog_system/pvt/products/GetProductAndSkuIds?_from=0&_to=1", baseURL(cred))

	_, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return err
	}

	return nil
}
