package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client/response"
)

func (c *shopifyClient) GetOrder(ctx context.Context, storeName, accessToken string, orderID string) (*domain.ShopifyOrder, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/orders/%s.json", storeName, orderID)

	var orderResp response.OrderResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&orderResp).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error al obtener orden de Shopify (código %d)", resp.StatusCode())
	}

	order := mappers.MapOrderResponseToShopifyOrder(orderResp.Order, nil, 0, "")

	return &order, nil
}
