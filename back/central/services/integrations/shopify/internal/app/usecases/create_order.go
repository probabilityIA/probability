package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"gorm.io/datatypes"
)

func (uc *SyncOrdersUseCase) CreateOrder(ctx context.Context, shopDomain string, order *domain.ShopifyOrder, rawPayload []byte) error {
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
	mapper.EnrichOrderWithDetails(probabilityOrder, rawPayload)

	// Agregar channel metadata con el payload original si está disponible
	if rawPayload != nil && len(rawPayload) > 0 {
		now := time.Now()
		probabilityOrder.ChannelMetadata = &domain.ProbabilityChannelMetadataDTO{
			ChannelSource: "shopify",
			RawData:       datatypes.JSON(rawPayload), // Convertir []byte a datatypes.JSON
			Version:       "1.0",
			ReceivedAt:    now,
			ProcessedAt:   &now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	return nil
}
