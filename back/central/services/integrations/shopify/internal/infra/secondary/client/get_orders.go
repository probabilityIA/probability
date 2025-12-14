package client

import (
	"context"
	"encoding/json"
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

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Convertir parámetros a query string
	queryParams := params.ToQueryString()
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Shopify-Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch orders, status: %d", resp.StatusCode)
	}

	// Decodificar la respuesta usando las estructuras tipadas
	var ordersResp response.OrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode orders response: %w", err)
	}

	// Mapear las órdenes tipadas a ShopifyOrder del dominio
	// Nota: businessID, integrationID e integrationType se completarán en el caso de uso
	// Por ahora retornamos órdenes sin estos campos
	orders := mappers.MapOrdersResponseToShopifyOrders(ordersResp.Orders, nil, 0, "")

	// Parse Link header for pagination
	linkHeader := resp.Header.Get("Link")
	nextPageURL := parseLinkHeader(linkHeader)

	return orders, nextPageURL, nil
}
