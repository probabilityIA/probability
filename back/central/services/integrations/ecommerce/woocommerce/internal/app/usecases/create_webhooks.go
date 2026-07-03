package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func (uc *wooCommerceUseCase) resolveStoreCreds(ctx context.Context, integrationID string) (storeURL, consumerKey, consumerSecret string, err error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return "", "", "", fmt.Errorf("getting integration: %w", err)
	}
	storeURL, err = extractString(integration.Config, "store_url")
	if err != nil || storeURL == "" {
		return "", "", "", fmt.Errorf("store_url not found in config")
	}
	consumerKey, err = uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return "", "", "", fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err = uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return "", "", "", fmt.Errorf("decrypting consumer_secret: %w", err)
	}
	return storeURL, consumerKey, consumerSecret, nil
}

func (uc *wooCommerceUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error) {
	storeURL, consumerKey, consumerSecret, err := uc.resolveStoreCreds(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	return uc.client.ListWebhooks(ctx, storeURL, consumerKey, consumerSecret)
}

func (uc *wooCommerceUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	storeURL, consumerKey, consumerSecret, err := uc.resolveStoreCreds(ctx, integrationID)
	if err != nil {
		return err
	}
	return uc.client.DeleteWebhook(ctx, storeURL, consumerKey, consumerSecret, webhookID)
}

func (uc *wooCommerceUseCase) CreateWebhooks(ctx context.Context, integrationID, baseURL, secret string) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}

	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil || storeURL == "" {
		return fmt.Errorf("store_url not found in config")
	}

	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	base := strings.TrimRight(baseURL, "/")
	base = strings.TrimSuffix(base, "/api/v1")
	deliveryURL := fmt.Sprintf("%s/api/v1/woocommerce/webhook?integration_id=%s", base, integrationID)

	topics := []string{"order.created", "order.updated"}
	ids := make([]int64, 0, len(topics))
	configured := true
	for _, topic := range topics {
		id, err := uc.client.CreateWebhook(ctx, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic)
		if err != nil {
			configured = false
			uc.logger.Error(ctx).Err(err).Str("topic", topic).Str("store_url", storeURL).Msg("Error creating WooCommerce webhook")
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return fmt.Errorf("could not create any WooCommerce webhook")
	}

	configUpdate := map[string]interface{}{
		"webhook_url":        deliveryURL,
		"webhook_configured": configured,
		"webhook_ids":        ids,
	}
	if err := uc.service.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return fmt.Errorf("updating integration config: %w", err)
	}
	return nil
}
