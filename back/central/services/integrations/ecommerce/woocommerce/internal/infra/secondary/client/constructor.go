package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// WooCommerceClient implementa domain.IWooCommerceClient usando la REST API de WooCommerce.
type WooCommerceClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para WooCommerce.
func New() domain.IWooCommerceClient {
	return &WooCommerceClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales llamando a GET /wp-json/wc/v3/system_status.
func (c *WooCommerceClient) TestConnection(ctx context.Context, storeURL, consumerKey, consumerSecret string) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/system_status", storeURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("la URL de la tienda no es valida: %w", err)
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con la tienda (verifica la URL y que el sitio este en linea): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("la tienda respondio 404: verifica que la URL apunte a un WordPress con WooCommerce y permalinks activos")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("la tienda respondio con estado inesperado %d", resp.StatusCode)
	}

	return nil
}
