package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) CreateOrder(ctx context.Context, shopDomain string, order *domain.ShopifyOrder, rawPayload []byte, isTest bool) error {
	if order == nil {
		return domain.ErrOrderPayloadNil
	}

	uc.log.Info(ctx).
		Str("shop_domain", shopDomain).
		Str("external_id", order.ExternalID).
		Msg("Processing orders/create webhook")

	integration, err := uc.integrationService.GetIntegrationByExternalID(ctx, shopDomain, domain.IntegrationTypeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("shop_domain", shopDomain).Msg("Failed to get integration by store domain")
		return fmt.Errorf("failed to get integration by store domain: %w", err)
	}

	if !integration.IsActive {
		uc.log.Warn(ctx).Uint("integration_id", integration.ID).Str("shop_domain", shopDomain).Msg("Integration is inactive, skipping webhook")
		return domain.ErrIntegrationInactive
	}

	order.BusinessID = integration.BusinessID
	order.IntegrationID = integration.ID
	order.IntegrationType = "shopify"

	probabilityOrder := mapper.MapShopifyOrderToProbability(order)

	// Marcar como orden de prueba si el header X-Probability-Testing estaba presente
	probabilityOrder.IsTest = isTest

	// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
	mapper.EnrichOrderWithDetails(probabilityOrder, rawPayload)

	// Agregar channel metadata con el payload original si está disponible
	if len(rawPayload) > 0 {
		now := time.Now()
		probabilityOrder.ChannelMetadata = &domain.ProbabilityChannelMetadataDTO{
			ChannelSource: "shopify",
			RawData:       rawPayload,
			Version:       "1.0",
			ReceivedAt:    now,
			ProcessedAt:   &now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
		uc.log.Error(ctx).Err(err).Str("external_id", order.ExternalID).Msg("Failed to publish order to queue")
		return fmt.Errorf("%w: %v", domain.ErrPublishFailed, err)
	}

	return nil
}
