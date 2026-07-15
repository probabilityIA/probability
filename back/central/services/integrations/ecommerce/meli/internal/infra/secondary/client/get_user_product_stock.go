package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type userProductStockResponse struct {
	Locations []struct {
		Type          string `json:"type"`
		StoreID       string `json:"store_id"`
		NetworkNodeID string `json:"network_node_id"`
		Quantity      int    `json:"quantity"`
	} `json:"locations"`
}

func (c *MeliClient) GetUserProductStock(ctx context.Context, accessToken, userProductID string) (*domain.UserProductStock, error) {
	endpoint := fmt.Sprintf("%s/user-products/%s/stock", c.baseURL, userProductID)

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
		return nil, domain.ErrItemNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: user-product stock status %d: %s", resp.StatusCode, string(body))
	}

	var parsed userProductStockResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("meli client: parsing user-product stock: %w", err)
	}

	stock := &domain.UserProductStock{Version: resp.Header.Get("x-version")}
	for _, l := range parsed.Locations {
		stock.Locations = append(stock.Locations, domain.StockLocation{
			Type:          l.Type,
			StoreID:       l.StoreID,
			NetworkNodeID: l.NetworkNodeID,
			Quantity:      l.Quantity,
		})
	}
	return stock, nil
}
