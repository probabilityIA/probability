package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) ProcessOrderPaid(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error {
	if order == nil {
		return domain.ErrOrderPayloadNil
	}

	uc.log.Info(ctx).
		Str("shop_domain", shopDomain).
		Str("external_id", order.ExternalID).
		Msg("Processing orders/paid webhook")

	integration, err := uc.integrationService.GetIntegrationByExternalID(ctx, shopDomain, domain.IntegrationTypeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("shop_domain", shopDomain).Msg("Failed to get integration by store domain")
		return fmt.Errorf("failed to get integration by store domain: %w", err)
	}

	order.BusinessID = integration.BusinessID
	order.IntegrationID = integration.ID
	order.IntegrationType = "shopify"

	probabilityOrder := mapper.MapShopifyOrderToProbability(order)

	// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
	// Estos detalles incluyen financial_status y fulfillment_status que se mapearán a PaymentStatusID y FulfillmentStatusID
	mapper.EnrichOrderWithDetails(probabilityOrder, order.RawData)

	if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
		uc.log.Error(ctx).Err(err).Str("external_id", order.ExternalID).Msg("Failed to publish order to queue")
		return fmt.Errorf("%w: %v", domain.ErrPublishFailed, err)
	}

	return nil
}
