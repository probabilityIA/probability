package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) GetOrder(ctx context.Context, integrationID string, orderID string) error {
	uc.log.Info(ctx).
		Str("integration_id", integrationID).
		Str("order_id", orderID).
		Msg("Getting single order from Shopify")

	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Failed to get integration")
		return fmt.Errorf("failed to get integration: %w", err)
	}

	config, err := utils.NormalizeConfig(integration.Config, integration.Name)
	if err != nil {
		return fmt.Errorf("invalid integration config: %w", err)
	}

	storeName, err := utils.ExtractStoreName(config, integration.Name)
	if err != nil {
		return fmt.Errorf("failed to extract store name: %w", err)
	}

	// En modo test, usar la URL de pruebas (base_url_test) en vez del dominio de Shopify
	storeName = utils.ResolveEffectiveStoreDomain(integration, storeName)

	accessToken, err := utils.GetAccessToken(ctx, uc.integrationService, integrationID)
	if err != nil {
		return err
	}

	shopifyOrder, err := uc.shopifyClient.GetOrder(ctx, storeName, accessToken, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order from Shopify: %w", err)
	}

	var order *domain.ShopifyOrder = shopifyOrder
	order.BusinessID = integration.BusinessID
	order.IntegrationID = integration.ID
	order.IntegrationType = "shopify"

	probabilityOrder := mapper.MapShopifyOrderToProbability(order)

	// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
	// Estos detalles incluyen financial_status y fulfillment_status que se mapearán a PaymentStatusID y FulfillmentStatusID
	mapper.EnrichOrderWithDetails(probabilityOrder, order.RawData)

	if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
		uc.log.Error(ctx).Err(err).Str("order_id", orderID).Msg("Failed to publish order to queue")
		return fmt.Errorf("%w: %v", domain.ErrPublishFailed, err)
	}

	return nil
}
