package usecaseorderscore

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

// CalculateAndUpdateOrderScore calcula el score de una orden y lo actualiza en la base de datos
// Este es el método principal que se activa mediante eventos
func (uc *UseCaseOrderScore) CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error {
	// 1. Obtener la orden de la base de datos
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// 2. Obtener datos necesarios para el cálculo del score
	// 2.1. Obtener historial de compra del cliente
	// ESTRATEGIA HÍBRIDA DE RECUPERACIÓN:
	// Si es orden local (IntegrationID==0), confiamos en la DB.
	// Si es integración (Shopify), normalmente confiamos en el dato mapeado.
	// PERO, si el dato mapeado es 0, podría ser un cliente nuevo (correcto) o una orden vieja sin dato (bug).
	// Solución: Si es 0, consultamos DB. Si DB > 1, es cliente recurrente (recuperamos historial). Si DB <= 1, es nuevo (mantenemos 0).
	if order.CustomerID != nil {
		// Validar si necesitamos consultar DB
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

	// 2.2. Obtener ShippingStreet2 desde Addresses si no está en el campo plano
	if order.Address2 == "" && len(order.Addresses) > 0 {
		for _, addr := range order.Addresses {
			if addr.Type == "shipping" && addr.Street2 != "" {
				order.Address2 = addr.Street2
				break
			}
		}
	}

	// 3. Calcular el score
	score, factors := uc.CalculateOrderScore(order)

	// 4. Actualizar la orden con el score y los factores negativos
	order.DeliveryProbability = &score

	// Serializar factors a JSON
	if len(factors) > 0 {
		factorsJSON, err := json.Marshal(factors)
		if err == nil {
			order.NegativeFactors = datatypes.JSON(factorsJSON)
		}
	} else {
		order.NegativeFactors = datatypes.JSON("[]")
	}

	// 5. Guardar los cambios
	if err := uc.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("failed to update order with score: %w", err)
	}

	return nil
}
