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

func (c *shopifyClient) GetOrders(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/orders.json", storeName)

	// Convertir parámetros a query string
	queryParams := params.ToQueryString()

	var ordersResp response.OrdersResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParams).
		SetResult(&ordersResp).
		Get(url)

	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, "", fmt.Errorf("error al obtener órdenes de Shopify (código %d)", resp.StatusCode())
	}

	// Mapear las órdenes tipadas a ShopifyOrder del dominio
	// Nota: businessID, integrationID e integrationType se completarán en el caso de uso
	// Por ahora retornamos órdenes sin estos campos
	orders := mappers.MapOrdersResponseToShopifyOrders(ordersResp.Orders, nil, 0, "")

	// Parse Link header for pagination
	linkHeader := resp.Header().Get("Link")
	nextPageURL := parseLinkHeader(linkHeader)

	return orders, nextPageURL, nil
}
