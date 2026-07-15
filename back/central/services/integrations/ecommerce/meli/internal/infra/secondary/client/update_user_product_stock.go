package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type sellerWarehouseLocation struct {
	StoreID       string `json:"store_id"`
	NetworkNodeID string `json:"network_node_id"`
	Quantity      int    `json:"quantity"`
}

type sellerWarehousePayload struct {
	Locations []sellerWarehouseLocation `json:"locations"`
}

func (c *MeliClient) UpdateUserProductStock(ctx context.Context, accessToken, userProductID, version string, locations []domain.StockLocation) error {
	payload := sellerWarehousePayload{}
	for _, l := range locations {
		payload.Locations = append(payload.Locations, sellerWarehouseLocation{
			StoreID:       l.StoreID,
			NetworkNodeID: l.NetworkNodeID,
			Quantity:      l.Quantity,
		})
	}

	if err := c.putSellerWarehouse(ctx, accessToken, userProductID, version, payload); err != nil {
		if err == errStockVersionConflict {
			fresh, gerr := c.GetUserProductStock(ctx, accessToken, userProductID)
			if gerr != nil {
				return gerr
			}
			return c.putSellerWarehouse(ctx, accessToken, userProductID, fresh.Version, payload)
		}
		return err
	}
	return nil
}

var errStockVersionConflict = fmt.Errorf("meli client: stock version conflict")

func (c *MeliClient) putSellerWarehouse(ctx context.Context, accessToken, userProductID, version string, payload sellerWarehousePayload) error {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s/user-products/%s/stock/type/seller_warehouse", c.baseURL, userProductID)

	resp, respBody, err := c.do(ctx, func() (*http.Request, error) {
		req, berr := c.newAuthorizedRequestWithBody(ctx, http.MethodPut, endpoint, accessToken, bodyBytes)
		if berr != nil {
			return nil, berr
		}
		if version != "" {
			req.Header.Set("x-version", version)
		}
		return req, nil
	})
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusConflict {
		return errStockVersionConflict
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("meli client: seller_warehouse status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
