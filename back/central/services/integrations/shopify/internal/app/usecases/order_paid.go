package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) ProcessOrderPaid(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error {
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
	if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}
