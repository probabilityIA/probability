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
)

type wooProductResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	SKU           string `json:"sku"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	StockQuantity *int   `json:"stock_quantity"`
}

type wooVariationResponse struct {
	ID            int64  `json:"id"`
	SKU           string `json:"sku"`
	Price         string `json:"price"`
	StockQuantity *int   `json:"stock_quantity"`
}

func (c *WooCommerceClient) GetProducts(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]domain.WooProduct, error) {
	storeURL = strings.TrimRight(storeURL, "/")

	products := make([]domain.WooProduct, 0)
	page := 1
	for {
		endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products?per_page=100&page=%d&status=any", storeURL, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: creating request: %w", err)
		}
		req.SetBasicAuth(consumerKey, consumerSecret)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			resp.Body.Close()
			return nil, domain.ErrInvalidCredentials
		}
		if resp.StatusCode != http.StatusOK {
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: unexpected status %d listing products: %s", resp.StatusCode, string(raw))
		}

		var list []wooProductResponse
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: decoding products response: %w", err)
		}
		resp.Body.Close()

		if len(list) == 0 {
			break
		}

		for _, p := range list {
			if p.Type == "variable" {
				variations, verr := c.getProductVariations(ctx, storeURL, consumerKey, consumerSecret, p.ID)
				if verr != nil {
					return nil, verr
				}
				parentID := strconv.FormatInt(p.ID, 10)
				for _, v := range variations {
					vprice := 0.0
					if v.Price != "" {
						vprice, _ = strconv.ParseFloat(v.Price, 64)
					}
					vstock := 0
					if v.StockQuantity != nil {
						vstock = *v.StockQuantity
					}
					products = append(products, domain.WooProduct{
						ID:            strconv.FormatInt(v.ID, 10),
						ParentID:      parentID,
						SKU:           v.SKU,
						Name:          p.Name,
						Price:         vprice,
						StockQuantity: vstock,
					})
				}
				continue
			}

			price := 0.0
			if p.Price != "" {
				price, _ = strconv.ParseFloat(p.Price, 64)
			}
			stock := 0
			if p.StockQuantity != nil {
				stock = *p.StockQuantity
			}
			products = append(products, domain.WooProduct{
				ID:            strconv.FormatInt(p.ID, 10),
				SKU:           p.SKU,
				Name:          p.Name,
				Price:         price,
				StockQuantity: stock,
			})
		}

		if len(list) < 100 {
			break
		}
		page++
	}

	return products, nil
}

func (c *WooCommerceClient) getProductVariations(ctx context.Context, storeURL, consumerKey, consumerSecret string, parentID int64) ([]wooVariationResponse, error) {
	all := make([]wooVariationResponse, 0)
	page := 1
	for {
		endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products/%d/variations?per_page=100&page=%d", storeURL, parentID, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: creating variations request: %w", err)
		}
		req.SetBasicAuth(consumerKey, consumerSecret)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: variations request failed: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: unexpected status %d listing variations of %d: %s", resp.StatusCode, parentID, string(raw))
		}

		var list []wooVariationResponse
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: decoding variations response: %w", err)
		}
		resp.Body.Close()

		if len(list) == 0 {
			break
		}
		all = append(all, list...)
		if len(list) < 100 {
			break
		}
		page++
	}
	return all, nil
}
