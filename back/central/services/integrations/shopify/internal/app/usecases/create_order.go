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

	// Procesar y extraer datos adicionales del JSON original
	if rawPayload != nil && len(rawPayload) > 0 {
		// Extraer FinancialDetails
		if financialDetails, err := mapper.ExtractFinancialDetails(rawPayload); err == nil {
			probabilityOrder.FinancialDetails = financialDetails
		}

		// Extraer ShippingDetails
		if shippingDetails, err := mapper.ExtractShippingDetails(rawPayload); err == nil {
			probabilityOrder.ShippingDetails = shippingDetails
		}

		// Extraer PaymentDetails
		if paymentDetails, err := mapper.ExtractPaymentDetails(rawPayload); err == nil {
			probabilityOrder.PaymentDetails = paymentDetails
		}

		// Extraer FulfillmentDetails
		if fulfillmentDetails, err := mapper.ExtractFulfillmentDetails(rawPayload); err == nil {
			probabilityOrder.FulfillmentDetails = fulfillmentDetails
		}

		// Agregar channel metadata con el payload original
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
