package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

type shopifyClient struct {
	httpClient *http.Client
}

func New() domain.ShopifyClient {
	return &shopifyClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *shopifyClient) ValidateToken(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/shop.json", storeName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, nil, err
	}

	req.Header.Set("X-Shopify-Access-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return true, nil, nil // Valid connection but failed to parse body
		}
		return true, result, nil
	}

	return false, nil, fmt.Errorf("shopify api returned status: %d", resp.StatusCode)
}

func (c *shopifyClient) FetchOrders(ctx context.Context, storeName, accessToken string, params map[string]string) ([]map[string]interface{}, string, error) {
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}

	url := fmt.Sprintf("https://%s/admin/api/2024-10/orders.json", storeName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	// Default params if not present
	if q.Get("status") == "" {
		q.Set("status", "any")
	}
	if q.Get("limit") == "" {
		q.Set("limit", "250")
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

	var result struct {
		Orders []map[string]interface{} `json:"orders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	// Parse Link header for pagination
	linkHeader := resp.Header.Get("Link")
	nextPageURL := parseLinkHeader(linkHeader)

	return result.Orders, nextPageURL, nil
}

func parseLinkHeader(header string) string {
	if header == "" {
		return ""
	}
	links := strings.Split(header, ",")
	for _, link := range links {
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}
		if strings.Contains(parts[1], `rel="next"`) {
			url := strings.Trim(parts[0], " <>")
			return url
		}
	}
	return ""
}
