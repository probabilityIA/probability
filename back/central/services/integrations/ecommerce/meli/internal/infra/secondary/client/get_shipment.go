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

// GetShipmentDetail obtiene los detalles completos de un envío.
// GET https://api.mercadolibre.com/shipments/{shipment_id}
// La orden solo trae shipping.id; esta llamada obtiene dirección, estado y opción de envío.
func (c *MeliClient) GetShipmentDetail(ctx context.Context, accessToken string, shipmentID int64) (*domain.MeliShippingDetail, error) {
	endpoint := fmt.Sprintf("%s/shipments/%d", c.baseURL, shipmentID)

	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrTokenExpired
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, domain.ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("meli client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var shippingResp response.MeliShippingDetailResponse
	if err := json.Unmarshal(body, &shippingResp); err != nil {
		return nil, fmt.Errorf("meli client: parsing shipping response: %w", err)
	}

	detail := shippingResp.ToDomain()
	return &detail, nil
}
