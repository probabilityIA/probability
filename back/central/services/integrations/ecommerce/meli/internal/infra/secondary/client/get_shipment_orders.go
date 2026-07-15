package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (c *MeliClient) GetShipmentOrderIDs(ctx context.Context, accessToken string, shipmentID int64) ([]int64, error) {
	endpoint := fmt.Sprintf("%s/shipments/%d/items", c.baseURL, shipmentID)

	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: shipment items status %d: %s", resp.StatusCode, string(body))
	}

	var items []struct {
		OrderID int64 `json:"order_id"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("meli client: parsing shipment items: %w", err)
	}

	seen := make(map[int64]bool)
	var ids []int64
	for _, it := range items {
		if it.OrderID > 0 && !seen[it.OrderID] {
			seen[it.OrderID] = true
			ids = append(ids, it.OrderID)
		}
	}
	return ids, nil
}
