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

type wooWebhookResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Topic       string `json:"topic"`
	DeliveryURL string `json:"delivery_url"`
	DateCreated string `json:"date_created"`
}

func (c *WooCommerceClient) ListWebhooks(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]domain.WebhookItem, error) {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/webhooks?per_page=100", storeURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("woocommerce client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("woocommerce client: unexpected status %d listing webhooks: %s", resp.StatusCode, string(raw))
	}

	var list []wooWebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("woocommerce client: decoding webhooks response: %w", err)
	}

	items := make([]domain.WebhookItem, 0, len(list))
	for _, w := range list {
		if !strings.Contains(w.DeliveryURL, "/woocommerce/webhook") {
			continue
		}
		items = append(items, domain.WebhookItem{
			ID:        strconv.FormatInt(w.ID, 10),
			Address:   w.DeliveryURL,
			Topic:     w.Topic,
			Format:    "json",
			CreatedAt: w.DateCreated,
		})
	}
	return items, nil
}

func (c *WooCommerceClient) DeleteWebhook(ctx context.Context, storeURL, consumerKey, consumerSecret, webhookID string) error {
	storeURL = strings.TrimRight(storeURL, "/")
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/webhooks/%s?force=true", storeURL, webhookID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("woocommerce client: creating request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("woocommerce client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("woocommerce client: unexpected status %d deleting webhook: %s", resp.StatusCode, string(raw))
	}
	return nil
}
