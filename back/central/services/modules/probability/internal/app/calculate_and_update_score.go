package app

import (
	"context"
	"encoding/json"
	"fmt"
)

func (uc *UseCaseScore) CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error {
	// 1. Get order for scoring
	order, err := uc.repo.GetOrderForScoring(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order for scoring: %w", err)
	}

	// 2. Hybrid customer history recovery
	// ESTRATEGIA HIBRIDA DE RECUPERACION:
	// Si es orden local (IntegrationID==0), confiamos en la DB.
	// Si es integracion (Shopify), normalmente confiamos en el dato mapeado.
	// PERO, si el dato mapeado es 0, podria ser un cliente nuevo (correcto) o una orden vieja sin dato (bug).
	// Solucion: Si es 0, consultamos DB. Si DB > 1, es cliente recurrente (recuperamos historial). Si DB <= 1, es nuevo (mantenemos 0).
	if order.CustomerID != nil {
		shouldCheckDB := order.IntegrationID == 0 || order.CustomerOrderCount == 0

		if shouldCheckDB {
			count, err := uc.repo.CountOrdersByClientID(ctx, *order.CustomerID)
			if err == nil {
				dbCount := int(count)

				if order.IntegrationID == 0 {
					// Para locales, DB es la verdad absoluta
					order.CustomerOrderCount = dbCount
				} else {
					// Para integraciones con count 0
					if dbCount > 1 {
						// Recuperamos historial perdido
						order.CustomerOrderCount = dbCount
					}
					// Si dbCount <= 1: es realmente nuevo, mantenemos 0 para que penalice
				}
			}
		}
	}

	// 3. Address2 fallback from Addresses
	if order.Address2 == "" && len(order.Addresses) > 0 {
		for _, addr := range order.Addresses {
			if addr.Type == "shipping" && addr.Street2 != "" {
				order.Address2 = addr.Street2
				break
			}
		}
	}

	// 4. Fetch customer order history (for purchase history, logistics, and payment risk categories)
	if order.CustomerID != nil {
		customerHistory, err := uc.repo.GetCustomerOrderHistory(ctx, *order.CustomerID, orderID)
		if err != nil {
			uc.log.Warn(ctx).Err(err).Str("order_id", orderID).Msg("Could not fetch customer order history, scoring will use defaults")
		} else {
			order.CustomerHistory = customerHistory
		}
	}

	// 5. Fetch customer delivery history (for logistics category)
	if order.CustomerID != nil {
		deliveryHistory, err := uc.repo.GetCustomerDeliveryHistory(ctx, *order.CustomerID)
		if err != nil {
			uc.log.Warn(ctx).Err(err).Str("order_id", orderID).Msg("Could not fetch customer delivery history, scoring will use defaults")
		} else {
			order.DeliveryHistory = deliveryHistory
		}
	}

	// 6. Fetch order item count (for order characteristics category)
	itemCount, err := uc.repo.GetOrderItemCount(ctx, orderID)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Str("order_id", orderID).Msg("Could not fetch order item count, scoring will use default 0")
	} else {
		order.OrderItemCount = itemCount
	}

	// 7. Fetch payment method category (for payment risk category)
	if order.PaymentMethodID > 0 {
		_, err := uc.repo.GetPaymentMethodCategory(ctx, order.PaymentMethodID)
		if err != nil {
			uc.log.Warn(ctx).Err(err).Str("order_id", orderID).Msg("Could not fetch payment method category")
		}
		// The category can be used in future scoring refinements
	}

	// 8. Calculate score (returns score, factors list, and full breakdown)
	score, factors, breakdown := uc.CalculateOrderScore(order)

	// 9. Serialize factors
	var factorsJSON []byte
	if len(factors) > 0 {
		factorsJSON, _ = json.Marshal(factors)
	} else {
		factorsJSON = []byte("[]")
	}

	// 10. Serialize breakdown
	var breakdownJSON []byte
	if breakdown != nil {
		breakdownJSON, _ = json.Marshal(breakdown)
	}

	// 11. Update order in DB (score + factors + breakdown)
	if err := uc.repo.UpdateOrderScore(ctx, orderID, score, factorsJSON, breakdownJSON); err != nil {
		return fmt.Errorf("failed to update order score: %w", err)
	}

	// 12. Publish score_calculated event
	businessID := uint(0)
	if order.BusinessID != nil {
		businessID = *order.BusinessID
	}
	if err := uc.publisher.PublishScoreCalculated(ctx, orderID, order.OrderNumber, businessID, order.IntegrationID); err != nil {
		uc.log.Error(ctx).Err(err).Str("order_id", orderID).Msg("Error publicando evento order.score_calculated")
	} else {
		uc.log.Info(ctx).Str("order_id", orderID).Str("order_number", order.OrderNumber).
			Msg("Score calculado y evento publicado exitosamente")
	}

	return nil
}
