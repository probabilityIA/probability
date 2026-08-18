package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		OrderID flexibleInt64 `json:"order_id"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("meli client: parsing shipment items: %w", err)
	}

	seen := make(map[int64]bool)
	var ids []int64
	for _, it := range items {
		id := int64(it.OrderID)
		if id > 0 && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids, nil
}

type flexibleInt64 int64

func (f *flexibleInt64) UnmarshalJSON(data []byte) error {
	trimmed := strings.Trim(string(data), "\"")
	if trimmed == "" || trimmed == "null" {
		*f = 0
		return nil
	}
	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return err
	}
	*f = flexibleInt64(value)
	return nil
}
