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
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Price        string `json:"price"`
	StockQuantity *int  `json:"stock_quantity"`
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
