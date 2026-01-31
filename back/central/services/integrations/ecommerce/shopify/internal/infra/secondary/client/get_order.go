package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

func (c *shopifyClient) GetOrder(ctx context.Context, storeName, accessToken string, orderID string) (*domain.ShopifyOrder, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/orders/%s.json", storeName, orderID)

	// 1. Perform Request
	var orderResp response.OrderResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error al obtener orden de Shopify (código %d)", resp.StatusCode())
	}

	// 2. Unmarshal into typed response
	if err := json.Unmarshal(resp.Body(), &orderResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling order response: %w", err)
	}

	// 2. Unmarshal into raw response to get original JSON
	var rawResp struct {
		Order json.RawMessage `json:"order"`
	}
	var rawOrder []byte
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		fmt.Printf("Warning: failed to unmarshal raw order: %v\n", err)
	} else {
		rawOrder = []byte(rawResp.Order)
	}

	order := mappers.MapOrderResponseToShopifyOrder(orderResp.Order, rawOrder, nil, 0, "")

	return &order, nil
}
