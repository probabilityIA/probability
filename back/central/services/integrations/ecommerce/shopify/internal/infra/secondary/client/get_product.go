package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (c *shopifyClient) GetProduct(ctx context.Context, storeName, accessToken, productID string) (*domain.ShopifyProduct, error) {
	url := buildURL(storeName, fmt.Sprintf("/admin/api/2024-10/products/%s.json", productID))

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("producto de Shopify no encontrado")
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error al obtener producto de Shopify (codigo %d)", resp.StatusCode())
	}

	var parsed struct {
		Product struct {
			ID       int64 `json:"id"`
			Variants []struct {
				ID              int64  `json:"id"`
				SKU             string `json:"sku"`
				InventoryItemID int64  `json:"inventory_item_id"`
			} `json:"variants"`
		} `json:"product"`
	}
	if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
		return nil, fmt.Errorf("error unmarshalling product: %w", err)
	}

	product := &domain.ShopifyProduct{ID: parsed.Product.ID}
	for _, v := range parsed.Product.Variants {
		product.Variants = append(product.Variants, domain.ShopifyVariant{
			ID:              v.ID,
			SKU:             v.SKU,
			InventoryItemID: v.InventoryItemID,
		})
	}
	return product, nil
}
