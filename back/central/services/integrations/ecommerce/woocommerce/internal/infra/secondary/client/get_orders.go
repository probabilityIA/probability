package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/client/response"
)

// GetOrders obtiene órdenes paginadas de WooCommerce REST API v3.
// Retorna las órdenes tipadas, los bytes crudos por orden (para ChannelMetadata.RawData), y error.
func (c *WooCommerceClient) GetOrders(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error) {
	storeURL = strings.TrimRight(storeURL, "/")

	queryStr := ""
	if params != nil {
		queryStr = params.ToQueryString()
	}

	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/orders", storeURL)
	if queryStr != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryStr)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: creating request: %w", err)
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("woocommerce client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: reading response: %w", err)
	}

	// Parsear como array de JSON crudos para preservar los bytes originales por orden
	var rawOrders []json.RawMessage
	if err := json.Unmarshal(body, &rawOrders); err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: parsing response: %w", err)
	}

	// Deserializar cada orden
	orders := make([]domain.WooCommerceOrder, 0, len(rawOrders))
	rawBytes := make([][]byte, 0, len(rawOrders))
	for _, raw := range rawOrders {
		var orderResp response.WooOrderResponse
		if err := json.Unmarshal(raw, &orderResp); err != nil {
			continue // skip malformed orders
		}
		orders = append(orders, orderResp.ToDomain())
		rawBytes = append(rawBytes, []byte(raw))
	}

	// Extraer headers de paginación
	total := parseHeaderInt(resp.Header.Get("X-WP-Total"))
	totalPages := parseHeaderInt(resp.Header.Get("X-WP-TotalPages"))

	return &domain.GetOrdersResult{
		Orders:     orders,
		Total:      total,
		TotalPages: totalPages,
	}, rawBytes, nil
}

// GetOrder obtiene una orden específica por ID.
func (c *WooCommerceClient) GetOrder(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*domain.WooCommerceOrder, []byte, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/orders/%d", storeURL, orderID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: creating request: %w", err)
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, domain.ErrNoOrdersFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("woocommerce client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: reading response: %w", err)
	}

	var orderResp response.WooOrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, nil, fmt.Errorf("woocommerce client: parsing response: %w", err)
	}

	order := orderResp.ToDomain()
	return &order, body, nil
}

func parseHeaderInt(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}
