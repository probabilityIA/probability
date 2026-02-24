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

// GetOrderByID obtiene el detalle completo de una orden de VTEX.
// GET https://{accountName}.vtexcommercestable.com.br/api/oms/pvt/orders/{orderId}
// Retorna la orden tipada, los bytes crudos (para ChannelMetadata.RawData), y error.
func (c *VTEXClient) GetOrderByID(ctx context.Context, storeURL, apiKey, apiToken string, orderID string) (*domain.VTEXOrder, []byte, error) {
	storeURL = normalizeStoreURL(storeURL)
	endpoint := fmt.Sprintf("%s/api/oms/pvt/orders/%s", storeURL, orderID)

	req, err := c.newVTEXRequest(ctx, http.MethodGet, endpoint, apiKey, apiToken)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("vtex client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, domain.ErrOrderNotFound
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, nil, domain.ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("vtex client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("vtex client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var orderResp response.VTEXOrderDetailResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, nil, fmt.Errorf("vtex client: parsing response: %w", err)
	}

	order := orderResp.ToDomain()
	return &order, body, nil
}
