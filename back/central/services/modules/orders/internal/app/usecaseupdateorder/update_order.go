package usecaseupdateorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// UpdateOrder actualiza una orden existente con los datos del DTO
func (uc *UseCaseUpdateOrder) UpdateOrder(ctx context.Context, existingOrder *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error) {
	// Guardar el estado anterior antes de actualizar (para detectar cambios de estado)
	previousStatus := existingOrder.Status

	// 1. Validar y actualizar todos los campos de la orden
	hasChanges := uc.updateOrderFields(ctx, existingOrder, dto)

	// 2. Si no hay cambios, retornar sin actualizar
	if !hasChanges {
		return uc.mapOrderToResponse(existingOrder), nil
	}

	// 3. Persistir los cambios
	if err := uc.repo.UpdateOrder(ctx, existingOrder); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	// 4. Publicar todos los eventos (integration sync, SSE, RabbitMQ, score)
	uc.publishUpdateEvents(ctx, existingOrder, previousStatus, dto.IsManualOrder)

	// 5. Retornar la respuesta actualizada
	return uc.mapOrderToResponse(existingOrder), nil
}

// updateOrderFields actualiza todos los campos de la orden y retorna si hubo cambios
func (uc *UseCaseUpdateOrder) updateOrderFields(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	hasChanges := false

	// Actualizar estados
	if uc.updateOrderStatuses(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar información financiera
	if uc.updateFinancialFields(order, dto) {
		hasChanges = true
	}

	// Actualizar información del cliente
	if uc.updateCustomerFields(order, dto) {
		hasChanges = true
	}

	// Actualizar dirección de envío
	if uc.updateShippingFields(order, dto) {
		hasChanges = true
	}

	// Actualizar información de pago
	if uc.updatePaymentFields(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar información de fulfillment
	if uc.updateFulfillmentFields(order, dto) {
		hasChanges = true
	}

	// Actualizar campos adicionales
	if uc.updateAdditionalFields(order, dto) {
		hasChanges = true
	}

	// Actualizar datos estructurados (JSONB)
	if uc.updateStructuredData(order, dto) {
		hasChanges = true
	}

	return hasChanges
}
