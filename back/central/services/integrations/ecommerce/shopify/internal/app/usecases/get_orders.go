package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
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

		// #region agent log
		if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "A",
				"location":     "get_orders.go:29",
				"message":      "GetOrders - Orders fetched from Shopify API",
				"data": map[string]interface{}{
					"orders_count":   len(orders),
					"next_url":       fetchedNextURL,
					"integration_id": integration.ID,
				},
				"timestamp": time.Now().UnixMilli(),
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion

		if integration.BusinessID == nil {
			err := fmt.Errorf("integration %d has no BusinessID assigned - cannot process orders", integration.ID)
			fmt.Printf("[GetOrders] ERR: %v\n", err)
			return err
		}

		publishedCount := 0
		publishErrorCount := 0
		for _, order := range orders {
			order.IntegrationID = integration.ID
			order.IntegrationType = "shopify"
			order.BusinessID = integration.BusinessID

			fmt.Printf("[GetOrders] Processing order ID: %s\n", order.ExternalID)
			probabilityOrder := mapper.MapShopifyOrderToProbability(&order)

			// Enriquecer la orden con detalles extraídos del JSON original (PaymentDetails, FulfillmentDetails, etc.)
			// Estos detalles incluyen financial_status y fulfillment_status que se mapearán a PaymentStatusID y FulfillmentStatusID
			mapper.EnrichOrderWithDetails(probabilityOrder, order.RawData)

			// #region agent log
			if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				logData, _ := json.Marshal(map[string]interface{}{
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "B",
					"location":     "get_orders.go:49",
					"message":      "GetOrders - Attempting to publish order to queue",
					"data": map[string]interface{}{
						"external_id":    order.ExternalID,
						"order_number":   probabilityOrder.OrderNumber,
						"integration_id": integration.ID,
					},
					"timestamp": time.Now().UnixMilli(),
				})
				f.WriteString(string(logData) + "\n")
				f.Close()
			}
			// #endregion

			if err := uc.orderPublisher.Publish(ctx, probabilityOrder); err != nil {
				fmt.Printf("[GetOrders] Error publishing order: %v. \n", err)
				// #region agent log
				if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					logData, _ := json.Marshal(map[string]interface{}{
						"sessionId":    "debug-session",
						"runId":        "run1",
						"hypothesisId": "B",
						"location":     "get_orders.go:52",
						"message":      "GetOrders - ERROR publishing order to queue",
						"data": map[string]interface{}{
							"external_id":    order.ExternalID,
							"error":          err.Error(),
							"integration_id": integration.ID,
						},
						"timestamp": time.Now().UnixMilli(),
					})
					f.WriteString(string(logData) + "\n")
					f.Close()
				}
				// #endregion
				publishErrorCount++
				// User requested NO fallback. Strict RabbitMQ usage.
				continue
			}
			// #region agent log
			if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				logData, _ := json.Marshal(map[string]interface{}{
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "B",
					"location":     "get_orders.go:54",
					"message":      "GetOrders - Order published successfully to queue",
					"data": map[string]interface{}{
						"external_id":    order.ExternalID,
						"order_number":   probabilityOrder.OrderNumber,
						"integration_id": integration.ID,
					},
					"timestamp": time.Now().UnixMilli(),
				})
				f.WriteString(string(logData) + "\n")
				f.Close()
			}
			// #endregion
			publishedCount++
			totalOrders++
		}

		// #region agent log
		if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logData, _ := json.Marshal(map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "B",
				"location":     "get_orders.go:56",
				"message":      "GetOrders - Page summary: published vs errors",
				"data": map[string]interface{}{
					"total_orders":    len(orders),
					"published_count": publishedCount,
					"error_count":     publishErrorCount,
					"integration_id":  integration.ID,
				},
				"timestamp": time.Now().UnixMilli(),
			})
			f.WriteString(string(logData) + "\n")
			f.Close()
		}
		// #endregion

		if fetchedNextURL == "" {
			break
		}
		nextURL = fetchedNextURL
	}

	fmt.Printf("[GetOrders] Completed: %d orders processed\n", totalOrders)
	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
			"location":     "get_orders.go:63",
			"message":      "GetOrders - Sync completed total summary",
			"data": map[string]interface{}{
				"total_orders_processed": totalOrders,
				"integration_id":         integration.ID,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion
	return nil
}
