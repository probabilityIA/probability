package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) ProcessOrderPartiallyFulfilled(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error {
	if order == nil {
		return fmt.Errorf("order payload is nil")
	}

	integration, err := uc.integrationService.GetIntegrationByStoreID(ctx, shopDomain)
	if err != nil {
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
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}
