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
	fmt.Printf("[CalculateAndUpdateOrderScore] ACTIVADO - Calculando score para orden ID: %s\n", orderID)

	// 1. Obtener la orden de la base de datos
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	fmt.Printf("[CalculateAndUpdateOrderScore] Orden %s obtenida - OrderNumber: %s, CustomerEmail: %s\n",
		orderID, order.OrderNumber, order.CustomerEmail)

	// 2. Obtener datos necesarios para el cálculo del score
	// 2.1. Obtener historial de compra del cliente
	if order.CustomerID != nil {
		count, err := uc.repo.CountOrdersByClientID(ctx, *order.CustomerID)
		if err == nil {
			order.CustomerOrderCount = int(count)
			fmt.Printf("[CalculateAndUpdateOrderScore] Orden %s - CustomerOrderCount: %d\n", orderID, order.CustomerOrderCount)
		}
	}

	// 2.2. Obtener ShippingStreet2 desde Addresses si no está en el campo plano
	if order.Address2 == "" && len(order.Addresses) > 0 {
		for _, addr := range order.Addresses {
			if addr.Type == "shipping" && addr.Street2 != "" {
				order.Address2 = addr.Street2
				fmt.Printf("[CalculateAndUpdateOrderScore] Orden %s - Address2 obtenido de Addresses: %s\n", orderID, order.Address2)
				break
			}
		}
	}

	// 3. Calcular el score
	score, factors := uc.CalculateOrderScore(order)
	fmt.Printf("[CalculateAndUpdateOrderScore] Orden %s - Score calculado: %.2f, Factors: %v\n", orderID, score, factors)

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

	fmt.Printf("[CalculateAndUpdateOrderScore] COMPLETADO - Orden %s actualizada con Score: %.2f\n", orderID, score)
	return nil
}
