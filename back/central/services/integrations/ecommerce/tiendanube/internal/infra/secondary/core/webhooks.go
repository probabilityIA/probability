package core

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func toInterfaceSlice(items []domain.WebhookItem) []interface{} {
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":         item.ID,
			"address":    item.Address,
			"topic":      item.Topic,
			"format":     item.Format,
			"created_at": item.CreatedAt,
			"updated_at": item.UpdatedAt,
		})
	}
	return result
}

func (t *TiendanubeCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	items, err := t.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	return toInterfaceSlice(items), nil
}

func (t *TiendanubeCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return t.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

func (t *TiendanubeCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	result, err := t.useCase.CreateWebhooks(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"webhook_url":       result.WebhookURL,
		"created_webhooks":  result.CreatedWebhooks,
		"failed_webhooks":   result.FailedWebhooks,
		"existing_webhooks": toInterfaceSlice(result.ExistingWebhooks),
	}, nil
}

func (t *TiendanubeCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	items, err := t.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	ours := stripQuery(usecases.WebhookDeliveryURL(baseURL, 0))
	matching := make([]domain.WebhookItem, 0, len(items))
	for _, item := range items {
		if stripQuery(item.Address) == ours {
			matching = append(matching, item)
		}
	}
	return toInterfaceSlice(matching), nil
}

func stripQuery(rawURL string) string {
	if idx := strings.Index(rawURL, "?"); idx >= 0 {
		return rawURL[:idx]
	}
	return rawURL
}
