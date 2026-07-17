package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func WebhookDeliveryURL(baseURL string, integrationID uint) string {
	return fmt.Sprintf("%s/api/v1/jumpseller/webhook?integration_id=%d", strings.TrimRight(baseURL, "/"), integrationID)
}

func (uc *jumpsellerUseCase) CreateWebhooks(ctx context.Context, integrationID, baseURL string) error {
	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		return fmt.Errorf("obteniendo informacion de la tienda Jumpseller: %w", err)
	}
	uc.persistStoreInfo(ctx, integrationID, integration, storeInfo)

	deliveryURL := WebhookDeliveryURL(baseURL, integration.ID)

	existing, err := uc.client.ListHooks(ctx, cred)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Msg("No se pudieron listar los webhooks existentes de Jumpseller")
		existing = nil
	}

	registered := make(map[string]bool, len(existing))
	for _, hook := range existing {
		if hook.Address == deliveryURL {
			registered[hook.Topic] = true
		}
	}

	var created int
	for _, event := range domain.WebhookOrderEvents {
		if registered[event] {
			continue
		}
		if _, err := uc.client.CreateHook(ctx, cred, event, deliveryURL); err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("event", event).
				Str("integration_id", integrationID).
				Msg("No se pudo registrar el webhook de Jumpseller (puede que ya exista)")
			continue
		}
		created++
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("delivery_url", deliveryURL).
		Str("store_code", storeInfo.Code).
		Int("created", created).
		Msg("Webhooks de Jumpseller registrados")

	return nil
}

func (uc *jumpsellerUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error) {
	_, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	return uc.client.ListHooks(ctx, cred)
}

func (uc *jumpsellerUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	_, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return err
	}
	return uc.client.DeleteHook(ctx, cred, webhookID)
}
