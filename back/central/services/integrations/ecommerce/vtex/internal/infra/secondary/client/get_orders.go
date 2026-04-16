package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client/response"
)

// GetOrders obtiene la lista de órdenes con paginación.
// GET https://{accountName}.vtexcommercestable.com.br/api/oms/pvt/orders?page={page}&per_page={perPage}
// Filtros soportados: f_creationDate, f_status, etc.
func (c *VTEXClient) GetOrders(ctx context.Context, storeURL, apiKey, apiToken string, page, perPage int, filters map[string]string) (*domain.VTEXOrderListResponse, error) {
	storeURL = normalizeStoreURL(storeURL)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	endpoint := fmt.Sprintf("%s/api/oms/pvt/orders?page=%d&per_page=%d", storeURL, page, perPage)

	for key, value := range filters {
		endpoint += fmt.Sprintf("&%s=%s", key, value)
	}

	req, err := c.newVTEXRequest(ctx, http.MethodGet, endpoint, apiKey, apiToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vtex client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, domain.ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vtex client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vtex client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var listResp response.VTEXOrderListAPIResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("vtex client: parsing response: %w", err)
	}

	result := listResp.ToDomain()
	return &result, nil
}
