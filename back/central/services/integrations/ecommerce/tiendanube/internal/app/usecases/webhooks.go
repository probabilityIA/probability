package usecases

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func WebhookDeliveryURL(baseURL string, integrationID uint) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	trimmed = strings.TrimSuffix(trimmed, "/api/v1")
	return fmt.Sprintf("%s/api/v1/tiendanube/webhook?integration_id=%d", trimmed, integrationID)
}

func isPubliclyReachable(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "https" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || host == "::1" || host == "host.docker.internal" {
		return false
	}
	if strings.HasPrefix(host, "127.") || strings.HasSuffix(host, ".local") {
		return false
	}
	return true
}

func (uc *tiendanubeUseCase) CreateWebhooks(ctx context.Context, integrationID, baseURL string) (*domain.CreateWebhooksResult, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, err
	}

	deliveryURL := WebhookDeliveryURL(baseURL, integration.ID)
	if !isPubliclyReachable(deliveryURL) {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Str("delivery_url", deliveryURL).
			Msg("La URL de webhooks no es publica: Tiendanube no podra entregar eventos")
		return nil, domain.ErrMissingWebhookBaseURL
	}

	existing, err := uc.client.ListWebhooks(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("listando los webhooks existentes de Tiendanube: %w", err)
	}

	registered := make(map[string]bool, len(existing))
	for _, hook := range existing {
		if hook.Address == deliveryURL {
			registered[hook.Topic] = true
		}
	}

	result := &domain.CreateWebhooksResult{
		WebhookURL:       deliveryURL,
		ExistingWebhooks: existing,
	}

	for _, event := range domain.WebhookEvents {
		if registered[event] {
			continue
		}
		if _, err := uc.client.CreateWebhook(ctx, cred, event, deliveryURL); err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("event", event).
				Msg("Tiendanube rechazo el registro de un webhook")
			result.FailedWebhooks = append(result.FailedWebhooks, event)
			continue
		}
		result.CreatedWebhooks = append(result.CreatedWebhooks, event)
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("delivery_url", deliveryURL).
		Int("created", len(result.CreatedWebhooks)).
		Int("failed", len(result.FailedWebhooks)).
		Msg("Webhooks de Tiendanube registrados")

	if len(result.FailedWebhooks) > 0 && len(result.CreatedWebhooks) == 0 {
		return result, domain.ErrWebhookCreationFailed
	}

	return result, nil
}

func (uc *tiendanubeUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, err
	}

	return uc.client.ListWebhooks(ctx, cred)
}

func (uc *tiendanubeUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return err
	}

	return uc.client.DeleteWebhook(ctx, cred, webhookID)
}
