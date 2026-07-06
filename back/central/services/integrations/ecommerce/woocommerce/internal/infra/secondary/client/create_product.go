package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func (c *WooCommerceClient) CreateProduct(ctx context.Context, storeURL, consumerKey, consumerSecret string, input domain.CreateProductInput) (string, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products", storeURL)

	payload := map[string]interface{}{
		"name":           input.Name,
		"type":           "simple",
		"sku":            input.SKU,
		"regular_price":  strconv.FormatFloat(input.Price, 'f', -1, 64),
		"description":    input.Description,
		"manage_stock":   input.ManageStock,
		"stock_quantity": input.StockQuantity,
		"status":         "publish",
	}

	if input.ImageURL != "" {
		payload["images"] = []map[string]interface{}{
			{"src": input.ImageURL},
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: marshaling product payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: create product request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", domain.ErrInvalidCredentials
	}

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("woocommerce client: estado inesperado %d al crear producto: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || result.ID == 0 {
		return "", fmt.Errorf("woocommerce client: respuesta invalida al crear producto: %s", string(raw))
	}

	return strconv.FormatInt(result.ID, 10), nil
}
