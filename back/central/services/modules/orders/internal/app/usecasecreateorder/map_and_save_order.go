package usecasecreateorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

// MapAndSaveOrder recibe una orden en formato canónico y la guarda en todas las tablas relacionadas
// Este es el punto de entrada principal para todas las integraciones después de mapear sus datos
func (uc *UseCaseCreateOrder) MapAndSaveOrder(ctx context.Context, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error) {
	// 0. Validar datos obligatorios de integración
	if dto.IntegrationID == 0 && !dto.IsManualOrder {
		return nil, errors.New("integration_id is required")
	}
	if dto.BusinessID == nil || *dto.BusinessID == 0 {
		return nil, errors.New("business_id is required")
	}

	// 1. Verificar si existe una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, dto.ExternalID, dto.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		existingOrder, err := uc.repo.GetOrderByExternalID(ctx, dto.ExternalID, dto.IntegrationID)
		if err != nil {
			return nil, fmt.Errorf("error getting existing order: %w", err)
		}
		return uc.updateUseCase.UpdateOrder(ctx, existingOrder, dto)
	}

	// 1.5. Validar/Crear Cliente
	client, err := uc.GetOrCreateCustomer(ctx, *dto.BusinessID, dto)
	if err != nil {
		return nil, fmt.Errorf("error processing customer: %w", err)
	}
	var clientID *uint
	if client != nil {
		clientID = &client.ID
	}

	// Para órdenes manuales, si el customer name está vacío pero tiene first/last name, componer
	if dto.CustomerName == "" && (dto.CustomerFirstName != "" || dto.CustomerLastName != "") {
		dto.CustomerName = fmt.Sprintf("%s %s", dto.CustomerFirstName, dto.CustomerLastName)
	}

	// 1.6. Mapear estados de la orden
	statusMapping := uc.mapOrderStatuses(ctx, dto)

	// 2. Crear la entidad de dominio ProbabilityOrder
	order := uc.buildOrderEntity(dto, clientID, statusMapping)

	// 2.1. Asignar PaymentMethodID desde el primer pago
	uc.assignPaymentMethodID(order, dto)

	// 2.1.1. Mantener IsPaid actualizado según PaymentStatusID
	uc.syncIsPaidFromPaymentStatus(ctx, order, statusMapping.PaymentStatusID)

	// 2.3. Popular campos JSONB y planos de dirección
	uc.populateOrderFields(order, dto)

	// 3. Guardar la orden principal (sin score por ahora, se calculará mediante evento)
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// 4-8. Guardar entidades relacionadas
	if err := uc.saveRelatedEntities(ctx, order, dto); err != nil {
		return nil, err
	}

	// 9. Publicar todos los eventos (integration sync, SSE, RabbitMQ, score)
	uc.publishOrderEvents(ctx, order, dto.IsManualOrder)

	// 10. Retornar la respuesta mapeada
	return uc.mapOrderToResponse(order), nil
}
