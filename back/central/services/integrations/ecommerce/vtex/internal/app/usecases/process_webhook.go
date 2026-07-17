package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) ProcessWebhook(ctx context.Context, payload *domain.VTEXWebhookPayload) error {
	if payload.OrderID == "" {
		uc.logger.Warn(ctx).Msg("Ignoring VTEX webhook with empty OrderId")
		return nil
	}

	if payload.IntegrationID == "" {
		uc.logger.Warn(ctx).
			Str("order_id", payload.OrderID).
			Msg("VTEX webhook sin integration_id, no se puede identificar la integracion")
		return domain.ErrIntegrationNotFound
	}

	uc.logger.Info(ctx).
		Str("order_id", payload.OrderID).
		Str("state", payload.State).
		Str("integration_id", payload.IntegrationID).
		Msg("Processing VTEX webhook")

	integration, err := uc.service.GetIntegrationByID(ctx, payload.IntegrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	if payload.Origin != nil && payload.Origin.Account != "" {
		expected := CleanAccountName(payload.Origin.Account)
		configured, cfgErr := extractString(integration.Config, "account_name")
		if cfgErr == nil && expected != "" && CleanAccountName(configured) != expected {
			uc.logger.Warn(ctx).
				Str("webhook_account", expected).
				Str("integration_account", CleanAccountName(configured)).
				Str("integration_id", payload.IntegrationID).
				Msg("La cuenta del webhook VTEX no coincide con la de la integracion, se descarta")
			return domain.ErrIntegrationNotFound
		}
	}

	cred, err := uc.resolveCredential(ctx, integration, payload.IntegrationID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Uint("integration_id", integration.ID).
			Msg("Failed to get VTEX credentials")
		return fmt.Errorf("getting credentials: %w", err)
	}

	order, rawJSON, err := uc.client.GetOrderByID(ctx, cred, payload.OrderID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_id", payload.OrderID).
			Msg("Failed to fetch VTEX order detail")
		return fmt.Errorf("fetching order: %w", err)
	}

	dto := mapper.MapVTEXOrderToProbability(order, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID

	if err := uc.publisher.Publish(ctx, dto); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_id", payload.OrderID).
			Msg("Failed to publish VTEX order to queue")
		return fmt.Errorf("publishing order: %w", err)
	}

	uc.logger.Info(ctx).
		Str("order_id", payload.OrderID).
		Str("status", order.Status).
		Uint("integration_id", integration.ID).
		Msg("VTEX order published successfully via webhook")

	return nil
}
