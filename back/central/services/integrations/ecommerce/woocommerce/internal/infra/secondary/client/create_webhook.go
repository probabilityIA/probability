package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func (c *WooCommerceClient) CreateWebhook(ctx context.Context, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic string) (int64, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/webhooks", storeURL)

	payload := map[string]interface{}{
		"name":         "Probability " + topic,
		"topic":        topic,
		"delivery_url": deliveryURL,
		"status":       "active",
	}
	if secret != "" {
		payload["secret"] = secret
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("woocommerce client: marshaling webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("woocommerce client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return 0, domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("woocommerce client: unexpected status %d creating webhook: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("woocommerce client: decoding webhook response: %w", err)
	}
	return result.ID, nil
}
