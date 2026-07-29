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
					ID      int64  `json:"id"`
					SKU     string `json:"sku"`
					Barcode string `json:"barcode"`
					ImageID *int64 `json:"image_id"`
				} `json:"variants"`
				Image *struct {
					Src string `json:"src"`
				} `json:"image"`
				Images []struct {
					ID  int64  `json:"id"`
					Src string `json:"src"`
				} `json:"images"`
			} `json:"products"`
		}
		if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
			return nil, fmt.Errorf("error unmarshalling products: %w", err)
		}

		for _, p := range parsed.Products {
			productID := strconv.FormatInt(p.ID, 10)
			mainImage := ""
			if p.Image != nil {
				mainImage = p.Image.Src
			} else if len(p.Images) > 0 {
				mainImage = p.Images[0].Src
			}
			byImageID := make(map[int64]string, len(p.Images))
			for _, img := range p.Images {
				byImageID[img.ID] = img.Src
			}
			for _, v := range p.Variants {
				image := mainImage
				if v.ImageID != nil {
					if src, ok := byImageID[*v.ImageID]; ok && src != "" {
						image = src
					}
				}
				products = append(products, domain.ShopifyProductForSync{
					ProductID: productID,
					VariantID: strconv.FormatInt(v.ID, 10),
					SKU:       v.SKU,
					Barcode:   v.Barcode,
					Name:      p.Title,
					ImageURL:  image,
				})
			}
		}

		url = parseLinkHeader(resp.Header().Get("Link"))
	}

	return products, nil
}
