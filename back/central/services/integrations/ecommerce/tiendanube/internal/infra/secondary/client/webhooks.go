package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/client/response"
)

func (c *TiendanubeClient) ListWebhooks(ctx context.Context, cred domain.Credential) ([]domain.WebhookItem, error) {
	query := url.Values{}
	query.Set("per_page", "200")

	raw, _, err := c.do(ctx, cred, http.MethodGet, "/webhooks", query, nil)
	if err != nil {
		return nil, err
	}

	var hooks []response.Webhook
	if err := json.Unmarshal(raw, &hooks); err != nil {
		return nil, fmt.Errorf("tiendanube client: parsing webhooks: %w", err)
	}

	items := make([]domain.WebhookItem, 0, len(hooks))
	for _, hook := range hooks {
		items = append(items, domain.WebhookItem{
			ID:        strconv.FormatInt(hook.ID, 10),
			Address:   hook.URL,
			Topic:     hook.Event,
			Format:    "json",
			CreatedAt: hook.CreatedAt,
			UpdatedAt: hook.UpdatedAt,
		})
	}
	return items, nil
}

func (c *TiendanubeClient) CreateWebhook(ctx context.Context, cred domain.Credential, event, webhookURL string) (string, error) {
	body := map[string]string{
		"event": event,
		"url":   webhookURL,
	}

	raw, _, err := c.do(ctx, cred, http.MethodPost, "/webhooks", nil, body)
	if err != nil {
		return "", err
	}

	var hook response.Webhook
	if err := json.Unmarshal(raw, &hook); err != nil {
		return "", fmt.Errorf("tiendanube client: parsing created webhook: %w", err)
	}

	return strconv.FormatInt(hook.ID, 10), nil
}

func (c *TiendanubeClient) DeleteWebhook(ctx context.Context, cred domain.Credential, webhookID string) error {
	_, _, err := c.do(ctx, cred, http.MethodDelete, "/webhooks/"+webhookID, nil, nil)
	return err
}
