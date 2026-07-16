package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (c *shopifyClient) ListProducts(ctx context.Context, storeName, accessToken string) ([]domain.ShopifyProductForSync, error) {
	products := make([]domain.ShopifyProductForSync, 0)
	url := buildURL(storeName, "/admin/api/2024-10/products.json?limit=250")

	for url != "" {
		resp, err := c.client.R().
			SetContext(ctx).
			SetHeader("X-Shopify-Access-Token", accessToken).
			SetHeader("Content-Type", "application/json").
			Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("error al listar productos de Shopify (codigo %d)", resp.StatusCode())
		}

		var parsed struct {
			Products []struct {
				ID       int64  `json:"id"`
				Title    string `json:"title"`
				Variants []struct {
					ID  int64  `json:"id"`
					SKU string `json:"sku"`
				} `json:"variants"`
			} `json:"products"`
		}
		if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
			return nil, fmt.Errorf("error unmarshalling products: %w", err)
		}

		for _, p := range parsed.Products {
			productID := strconv.FormatInt(p.ID, 10)
			for _, v := range p.Variants {
				products = append(products, domain.ShopifyProductForSync{
					ProductID: productID,
					SKU:       v.SKU,
					Name:      p.Title,
				})
			}
		}

		url = parseLinkHeader(resp.Header().Get("Link"))
	}

	return products, nil
}
