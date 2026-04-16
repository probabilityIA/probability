package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

// GetOrder obtiene una orden específica por ID.
// GET https://api.mercadolibre.com/orders/{order_id}
// Retorna la orden tipada, los bytes crudos (para ChannelMetadata.RawData), y error.
func (c *MeliClient) GetOrder(ctx context.Context, accessToken string, orderID int64) (*domain.MeliOrder, []byte, error) {
	endpoint := fmt.Sprintf("%s/orders/%d", c.baseURL, orderID)

	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil, domain.ErrTokenExpired
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, domain.ErrOrderNotFound
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, nil, domain.ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("meli client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("meli client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var orderResp response.MeliOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, nil, fmt.Errorf("meli client: parsing response: %w", err)
	}

	order := orderResp.ToDomain()
	return &order, body, nil
}
