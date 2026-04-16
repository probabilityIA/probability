package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) GetOrders(ctx context.Context, integration *domain.Integration, storeDomain, accessToken string, params *domain.GetOrdersParams) (int, error) {
	totalOrders := 0
	nextURL := ""

	for {
		var orders []domain.ShopifyOrder
		var fetchedNextURL string
		var err error

		if nextURL == "" {
			uc.log.Info(ctx).Msg("Fetching first page of orders")
			orders, fetchedNextURL, err = uc.shopifyClient.GetOrders(ctx, storeDomain, accessToken, params)
		} else {
			time.Sleep(500 * time.Millisecond)
			uc.log.Info(ctx).Str("next_url", nextURL).Msg("Fetching next page of orders")
			orders, fetchedNextURL, err = uc.shopifyClient.GetOrdersByURL(ctx, nextURL, accessToken)
		}

		if err != nil {
			return totalOrders, fmt.Errorf("error fetching orders: %w", err)
		}

		uc.log.Info(ctx).
			Int("fetched_count", len(orders)).
			Str("next_url", fetchedNextURL).
			Msg("Orders page fetched")

		if integration.BusinessID == nil {
			uc.log.Error(ctx).
				Uint("integration_id", integration.ID).
				Msg("Integration has no BusinessID assigned")
			return totalOrders, fmt.Errorf("integration %d: %w", integration.ID, domain.ErrBusinessIDMissing)
		}

		publishedCount := 0
		publishErrorCount := 0
		for _, order := range orders {
			order.IntegrationID = integration.ID
			order.IntegrationType = "shopify"
			order.BusinessID = integration.BusinessID

			probabilityOrder := mapper.MapShopifyOrderToProbability(&order)

			// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
			// Estos detalles incluyen financial_status y fulfillment_status que se mapearán a PaymentStatusID y FulfillmentStatusID
			mapper.EnrichOrderWithDetails(probabilityOrder, order.RawData)

			if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
				uc.log.Error(ctx).
					Err(err).
					Str("external_id", order.ExternalID).
					Msg("Error publishing order to queue")
				publishErrorCount++
				// User requested NO fallback. Strict RabbitMQ usage.
				continue
			}
			publishedCount++
			totalOrders++
		}

		if publishErrorCount > 0 {
			uc.log.Warn(ctx).
				Int("publish_errors", publishErrorCount).
				Int("published", publishedCount).
				Int("total_in_page", len(orders)).
				Msg("Some orders failed to publish in this page")
		}

		if fetchedNextURL == "" {
			break
		}
		nextURL = fetchedNextURL
	}

	uc.log.Info(ctx).
		Int("total_published", totalOrders).
		Msg("GetOrders completed")
	return totalOrders, nil
}
