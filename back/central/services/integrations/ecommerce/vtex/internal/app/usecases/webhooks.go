package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func WebhookDeliveryURL(baseURL string, integrationID uint) string {
	return fmt.Sprintf("%s/api/v1/vtex/webhook?integration_id=%d", strings.TrimRight(baseURL, "/"), integrationID)
}

func WebhookHookKey(integrationID uint) string {
	return fmt.Sprintf("vtex-order-webhook-%d", integrationID)
}

func sameHookURL(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func (uc *vtexUseCase) resolveForWebhook(ctx context.Context, integrationID string) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.Credential{}, domain.ErrIntegrationNotFound
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	return integration, cred, nil
}

func (uc *vtexUseCase) InspectWebhook(ctx context.Context, integrationID, baseURL string) (*domain.WebhookItem, error) {
	integration, cred, err := uc.resolveForWebhook(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	hook, err := uc.client.GetOrderHook(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("consultando el hook de VTEX: %w", err)
	}
	if hook == nil {
		return nil, nil
	}

	ours := WebhookDeliveryURL(baseURL, integration.ID)

	return &domain.WebhookItem{
		ID:       cred.AccountName,
		Address:  hook.URL,
		Statuses: hook.Statuses,
		IsOurs:   sameHookURL(hook.URL, ours),
	}, nil
}

func (uc *vtexUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error) {
	item, err := uc.InspectWebhook(ctx, integrationID, uc.webhookBaseURL)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return []domain.WebhookItem{}, nil
	}
	return []domain.WebhookItem{*item}, nil
}

func (uc *vtexUseCase) CreateWebhooks(ctx context.Context, integrationID, baseURL string, force bool) error {
	integration, cred, err := uc.resolveForWebhook(ctx, integrationID)
	if err != nil {
		return err
	}

	deliveryURL := WebhookDeliveryURL(baseURL, integration.ID)

	existing, err := uc.client.GetOrderHook(ctx, cred)
	if err != nil {
		return fmt.Errorf("consultando el hook actual de VTEX: %w", err)
	}

	if existing != nil && !sameHookURL(existing.URL, deliveryURL) && !force {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Str("account", cred.AccountName).
			Str("existing_hook", existing.URL).
			Msg("La cuenta VTEX ya tiene un hook de otra herramienta, no se sobreescribe sin confirmacion")
		return fmt.Errorf("%w: %s", domain.ErrForeignHookExists, existing.URL)
	}

	if existing != nil && sameHookURL(existing.URL, deliveryURL) && existing.HasKey {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("El hook de VTEX ya estaba registrado y apunta a Probability")
		return nil
	}

	if err := uc.client.SetOrderHook(ctx, cred, deliveryURL, WebhookHookKey(integration.ID)); err != nil {
		return fmt.Errorf("registrando el hook en VTEX: %w", err)
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("account", cred.AccountName).
		Str("url", deliveryURL).
		Bool("replaced_foreign", existing != nil && !sameHookURL(existing.URL, deliveryURL)).
		Msg("Hook de ordenes registrado en VTEX")

	return nil
}

func (uc *vtexUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	_, cred, err := uc.resolveForWebhook(ctx, integrationID)
	if err != nil {
		return err
	}

	if err := uc.client.DeleteOrderHook(ctx, cred); err != nil {
		return fmt.Errorf("eliminando el hook de VTEX: %w", err)
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("account", cred.AccountName).
		Msg("Hook de ordenes eliminado en VTEX")

	return nil
}
