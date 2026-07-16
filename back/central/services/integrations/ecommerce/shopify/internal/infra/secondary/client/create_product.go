package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (c *shopifyClient) CreateProduct(ctx context.Context, storeName, accessToken string, input domain.CreateProductInput) (string, error) {
	url := buildURL(storeName, "/admin/api/2024-10/products.json")

	payload := map[string]interface{}{
		"product": map[string]interface{}{
			"title":     input.Name,
			"body_html": input.Description,
			"status":    "active",
			"variants": []map[string]interface{}{
				{
					"sku":                  input.SKU,
					"price":                fmt.Sprintf("%.2f", input.Price),
					"inventory_management": "shopify",
				},
			},
		},
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusCreated && resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error al crear producto en Shopify (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	var parsed struct {
		Product struct {
			ID int64 `json:"id"`
		} `json:"product"`
	}
	if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
		return "", fmt.Errorf("error unmarshalling created product: %w", err)
	}

	return strconv.FormatInt(parsed.Product.ID, 10) + ":" + input.SKU, nil
}
