package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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
