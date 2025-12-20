package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) GetOrders(ctx context.Context, integration *domain.Integration, storeDomain, accessToken string, params *domain.GetOrdersParams) error {
	totalOrders := 0
	nextURL := ""

	for {
		if nextURL == "" {
			fmt.Println("[GetOrders] Fetching first page...")
		} else {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("[GetOrders] Fetching next page: %s\n", nextURL)
		}

		orders, fetchedNextURL, err := uc.shopifyClient.GetOrders(ctx, storeDomain, accessToken, params)
		if err != nil {
			return fmt.Errorf("error fetching orders: %w", err)
		}

		fmt.Printf("[GetOrders] Fetched %d orders. NextURL: %s\n", len(orders), fetchedNextURL)

		if integration.BusinessID == nil {
			err := fmt.Errorf("integration %d has no BusinessID assigned - cannot process orders", integration.ID)
			fmt.Printf("[GetOrders] ERR: %v\n", err)
			return err
		}

		for _, order := range orders {
			order.IntegrationID = integration.ID
			order.IntegrationType = "shopify"
			order.BusinessID = integration.BusinessID

			fmt.Printf("[GetOrders] Processing order ID: %s\n", order.ExternalID)
			probabilityOrder := mapper.MapShopifyOrderToProbability(&order)

			// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
			// Estos detalles incluyen financial_status y fulfillment_status que se mapearán a PaymentStatusID y FulfillmentStatusID
			mapper.EnrichOrderWithDetails(probabilityOrder, order.RawData)

			if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
				fmt.Printf("[GetOrders] Error publishing order: %v. \n", err)
				// User requested NO fallback. Strict RabbitMQ usage.
				continue
			}
			totalOrders++
		}

		if fetchedNextURL == "" {
			break
		}
		nextURL = fetchedNextURL
	}

	fmt.Printf("[GetOrders] Completed: %d orders processed\n", totalOrders)
	return nil
}
