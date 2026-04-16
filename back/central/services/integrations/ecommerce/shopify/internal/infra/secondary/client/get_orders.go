package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client/response"
)

func (c *shopifyClient) GetOrders(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
	url := buildURL(storeName, "/admin/api/2024-10/orders.json")

	// Convertir parámetros a query string
	queryParams := params.ToQueryString()

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParams).
		Get(url)

	if err != nil {
		return nil, "", err
	}

	return c.parseOrdersResponse(resp)
}

// GetOrdersByURL fetches orders using a full pagination URL (from Shopify Link header).
// Shopify cursor-based pagination requires using the exact URL from the Link header.
func (c *shopifyClient) GetOrdersByURL(ctx context.Context, nextPageURL, accessToken string) ([]domain.ShopifyOrder, string, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		Get(nextPageURL)

	if err != nil {
		return nil, "", err
	}

	return c.parseOrdersResponse(resp)
}

// parseOrdersResponse parses the HTTP response from Shopify orders endpoint.
func (c *shopifyClient) parseOrdersResponse(resp interface{ StatusCode() int; Body() []byte; Header() http.Header }) ([]domain.ShopifyOrder, string, error) {
	if resp.StatusCode() != http.StatusOK {
		return nil, "", fmt.Errorf("error al obtener órdenes de Shopify (código %d)", resp.StatusCode())
	}

	// 1. Unmarshal into typed response
	var ordersResp response.OrdersResponse
	if err := json.Unmarshal(resp.Body(), &ordersResp); err != nil {
		return nil, "", fmt.Errorf("error unmarshalling orders response: %w", err)
	}

	// 2. Unmarshal into raw response to get original JSONs
	var rawResp struct {
		Orders []json.RawMessage `json:"orders"`
	}
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		fmt.Printf("Warning: failed to unmarshal raw orders: %v\n", err)
	}

	// Convert []json.RawMessage to [][]byte
	rawOrdersBytes := make([][]byte, len(rawResp.Orders))
	for i, raw := range rawResp.Orders {
		rawOrdersBytes[i] = []byte(raw)
	}

	orders := mappers.MapOrdersResponseToShopifyOrders(ordersResp.Orders, rawOrdersBytes, nil, 0, "")

	// Parse Link header for pagination
	linkHeader := resp.Header().Get("Link")
	nextPageURL := parseLinkHeader(linkHeader)

	return orders, nextPageURL, nil
}
