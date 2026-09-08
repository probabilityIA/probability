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

func (c *WooCommerceClient) CreateVariableProduct(ctx context.Context, storeURL, consumerKey, consumerSecret string, input domain.CreateVariableProductInput) (string, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products", storeURL)

	attributes := make([]map[string]interface{}, 0, len(input.Attributes))
	for _, attr := range input.Attributes {
		attributes = append(attributes, map[string]interface{}{
			"name":      attr.Name,
			"visible":   true,
			"variation": true,
			"options":   attr.Options,
		})
	}

	payload := map[string]interface{}{
		"name":        input.Name,
		"type":        "variable",
		"description": input.Description,
		"status":      "publish",
		"attributes":  attributes,
	}

	if input.ImageURL != "" {
		payload["images"] = []map[string]interface{}{{"src": input.ImageURL}}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: marshaling variable product payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: create variable product request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", domain.ErrInvalidCredentials
	}

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("woocommerce client: estado inesperado %d al crear producto variable: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || result.ID == 0 {
		return "", fmt.Errorf("woocommerce client: respuesta invalida al crear producto variable: %s", string(raw))
	}

	return strconv.FormatInt(result.ID, 10), nil
}

func (c *WooCommerceClient) CreateProductVariation(ctx context.Context, storeURL, consumerKey, consumerSecret, parentID string, input domain.CreateVariationInput) (string, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products/%s/variations", storeURL, parentID)

	attributes := make([]map[string]interface{}, 0, len(input.Attributes))
	for name, value := range input.Attributes {
		attributes = append(attributes, map[string]interface{}{
			"name":   name,
			"option": value,
		})
	}

	payload := map[string]interface{}{
		"sku":            input.SKU,
		"regular_price":  strconv.FormatFloat(input.Price, 'f', -1, 64),
		"manage_stock":   input.ManageStock,
		"stock_quantity": input.StockQuantity,
		"attributes":     attributes,
	}

	if input.ImageURL != "" {
		payload["image"] = map[string]interface{}{"src": input.ImageURL}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: marshaling variation payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("woocommerce client: create variation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", domain.ErrInvalidCredentials
	}

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("woocommerce client: estado inesperado %d al crear variacion: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || result.ID == 0 {
		return "", fmt.Errorf("woocommerce client: respuesta invalida al crear variacion: %s", string(raw))
	}

	return strconv.FormatInt(result.ID, 10), nil
}
