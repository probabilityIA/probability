package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

// TiendanubeClient implementa domain.ITiendanubeClient usando la REST API de Tiendanube.
type TiendanubeClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para Tiendanube.
func New() domain.ITiendanubeClient {
	return &TiendanubeClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales llamando a GET {storeURL}/v1/store.
func (c *TiendanubeClient) TestConnection(ctx context.Context, storeURL, accessToken string) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/v1/store", storeURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("tiendanube client: creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tiendanube client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tiendanube client: unexpected status %d", resp.StatusCode)
	}

	return nil
}
