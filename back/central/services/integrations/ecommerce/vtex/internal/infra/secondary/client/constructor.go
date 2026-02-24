package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// VTEXClient implementa domain.IVTEXClient usando la REST API de VTEX.
type VTEXClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para VTEX.
func New() domain.IVTEXClient {
	return &VTEXClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// newVTEXRequest crea un request HTTP con los headers de autenticación de VTEX.
func (c *VTEXClient) newVTEXRequest(ctx context.Context, method, url, apiKey, apiToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vtex client: creating request: %w", err)
	}
	req.Header.Set("X-VTEX-API-AppKey", apiKey)
	req.Header.Set("X-VTEX-API-AppToken", apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// normalizeStoreURL normaliza la URL del store eliminando trailing slashes.
func normalizeStoreURL(storeURL string) string {
	return strings.TrimRight(storeURL, "/")
}

// TestConnection verifica las credenciales llamando a GET /api/catalog/pvt/product?_offset=0&_limit=1.
// VTEX usa headers X-VTEX-API-AppKey y X-VTEX-API-AppToken para autenticación.
func (c *VTEXClient) TestConnection(ctx context.Context, storeURL, apiKey, apiToken string) error {
	storeURL = normalizeStoreURL(storeURL)
	endpoint := fmt.Sprintf("%s/api/catalog/pvt/product?_offset=0&_limit=1", storeURL)

	req, err := c.newVTEXRequest(ctx, http.MethodGet, endpoint, apiKey, apiToken)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vtex client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vtex client: unexpected status %d", resp.StatusCode)
	}

	return nil
}
