package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento/internal/domain"
)

// MagentoClient implementa domain.IMagentoClient usando la REST API de Magento.
type MagentoClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para Magento.
func New() domain.IMagentoClient {
	return &MagentoClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales llamando a GET {storeURL}/rest/V1/store/storeConfigs.
func (c *MagentoClient) TestConnection(ctx context.Context, storeURL, accessToken string) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/rest/V1/store/storeConfigs", storeURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("magento client: creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("magento client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("magento client: unexpected status %d", resp.StatusCode)
	}

	return nil
}
