package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("la tienda respondio con una pagina web, no con la API REST: verifica que la REST API este habilitada y que la tienda NO este en modo 'Coming soon' o mantenimiento")
	}

	var apiError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Code != "" {
		if strings.Contains(apiError.Code, "authentication") || strings.Contains(apiError.Code, "cannot_view") || strings.Contains(apiError.Code, "unauthorized") {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("la tienda respondio con un error: %s", apiError.Message)
	}

	var systemStatus struct {
		Environment json.RawMessage `json:"environment"`
	}
	if err := json.Unmarshal(body, &systemStatus); err != nil || len(systemStatus.Environment) == 0 {
		return domain.ErrInvalidCredentials
	}

	return nil
}
