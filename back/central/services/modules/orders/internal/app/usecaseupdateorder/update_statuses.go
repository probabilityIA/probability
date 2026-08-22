package usecaseupdateorder

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/helpers/statusmapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateOrderStatuses actualiza los estados de la orden (OrderStatus, PaymentStatus, FulfillmentStatus)
func (uc *UseCaseUpdateOrder) updateOrderStatuses(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	hasChanges := false

	// Actualizar OrderStatus
	if uc.updateOrderStatus(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar PaymentStatus
	if uc.updatePaymentStatus(ctx, order, dto) {
		hasChanges = true
	}

	// Actualizar FulfillmentStatus
	if uc.updateFulfillmentStatus(ctx, order, dto) {
		hasChanges = true
	}

	return hasChanges
}

// channelStatusInboundEnabled indica si esta integración acepta que el canal
// mueva el estado de sus órdenes. Por defecto sí: apagarlo es una decisión
// explícita del negocio.
func (uc *UseCaseUpdateOrder) channelStatusInboundEnabled(ctx context.Context, order *entities.ProbabilityOrder) bool {
	if order.IntegrationID == 0 {
		return true
	}
	habilitado, err := uc.repo.IsChannelStatusInboundEnabled(ctx, order.IntegrationID)
	if err != nil {
		uc.logger.Warn().
			Err(err).
			Str("order_id", order.ID).
			Msg("No se pudo leer la configuración de entrada de estados; se acepta el estado del canal")
		return true
	}
	return habilitado
}

// updateOrderStatus actualiza el estado general de la orden a partir de lo que
// reporta la integración. No aplica las reglas de transición internas: las
// estrictas viven en PUT /orders/:id/status (usecaseupdatestatus).
func (uc *UseCaseUpdateOrder) updateOrderStatus(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	currentStatus := entities.OrderStatus(order.Status)
	targetStatus := entities.OrderStatus(dto.Status)

	// Cada integración decide si el canal puede mover el estado de sus órdenes.
	// Con la entrada apagada solo pasan cancelaciones y reembolsos.
	channelWins := uc.channelStatusInboundEnabled(ctx, order) || targetStatus.IsChannelOverride()

	if dto.Status != "" && order.Status != dto.Status {
		if !channelWins {
			uc.logger.Info().
				Str("order_id", order.ID).
				Str("current", order.Status).
				Str("channel_status", dto.Status).
				Str("integration_type", dto.IntegrationType).
				Msg("La integración tiene apagada la entrada de estados; se ignora el estado que llega del canal")
		} else {
			if !currentStatus.CanTransitionTo(targetStatus) && targetStatus != entities.OrderStatusCancelled {
				uc.logger.Warn().
					Str("order_id", order.ID).
					Str("from", order.Status).
					Str("to", dto.Status).
					Str("integration_type", dto.IntegrationType).
					Msg("Integración realizó salto de estado que no cumple flujo v2 — aceptado por ser fuente externa")
			}
			order.Status = dto.Status
			changed = true
		}
	}

	// Actualizar OriginalStatus y mapear StatusID
	if dto.OriginalStatus != "" && order.OriginalStatus != dto.OriginalStatus {
		order.OriginalStatus = dto.OriginalStatus
		changed = true

		// El estado del canal se registra siempre, pero el StatusID solo se
		// remapea cuando el canal aún puede mandar sobre esta orden.
		if channelWins {
			mappedStatusID := uc.mapOrderStatusID(ctx, dto)
			if order.StatusID == nil || (mappedStatusID != nil && *order.StatusID != *mappedStatusID) || (mappedStatusID == nil && order.StatusID != nil) {
				order.StatusID = mappedStatusID
				changed = true
			}
		}
	}

	return changed
}

// updatePaymentStatus actualiza el estado de pago
func (uc *UseCaseUpdateOrder) updatePaymentStatus(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	// Mapear PaymentStatusID desde el DTO
	mappedPaymentStatusID := uc.mapPaymentStatusID(ctx, dto)

	// Si se mapeó un nuevo estado y es diferente al actual, actualizarlo
	if mappedPaymentStatusID != nil {
		if order.PaymentStatusID == nil || *order.PaymentStatusID != *mappedPaymentStatusID {
			order.PaymentStatusID = mappedPaymentStatusID
			changed = true
		}
	}

	// Sincronizar IsPaid basado en PaymentStatusID
	oldIsPaid := order.IsPaid
	uc.syncIsPaidFromPaymentStatus(ctx, order, order.PaymentStatusID)
	if order.IsPaid != oldIsPaid {
		changed = true
	}

	return changed
}

// updateFulfillmentStatus actualiza el estado de fulfillment
func (uc *UseCaseUpdateOrder) updateFulfillmentStatus(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	// Mapear FulfillmentStatusID desde el DTO
	mappedFulfillmentStatusID := uc.mapFulfillmentStatusID(ctx, dto)

	// Si se mapeó un nuevo estado y es diferente al actual, actualizarlo
	if mappedFulfillmentStatusID != nil {
		if order.FulfillmentStatusID == nil || *order.FulfillmentStatusID != *mappedFulfillmentStatusID {
			order.FulfillmentStatusID = mappedFulfillmentStatusID
			changed = true
		}
	}

	return changed
}

// mapOrderStatusID mapea el estado general de la orden (replicado para independencia de módulos)
func (uc *UseCaseUpdateOrder) mapOrderStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	if dto.IntegrationType != "" {
		integrationTypeID := statusmapper.GetIntegrationTypeID(dto.IntegrationType)
		if integrationTypeID > 0 {
			if dto.Status != "" {
				mappedStatusID, err := uc.repo.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.Status)
				if err == nil && mappedStatusID != nil {
					return mappedStatusID
				}
			}

			if dto.OriginalStatus != "" {
				mappedStatusID, err := uc.repo.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.OriginalStatus)
				if err != nil {
					uc.logger.Warn().
						Uint("integration_type_id", integrationTypeID).
						Str("status", dto.Status).
						Str("original_status", dto.OriginalStatus).
						Err(err).
						Msg("Error al buscar mapeo de estado, continuando sin status_id")
				}
				if mappedStatusID != nil {
					return mappedStatusID
				}
			}
		}
	}

	if dto.Status != "" {
		statusID, err := uc.repo.GetOrderStatusIDByCode(ctx, dto.Status)
		if err == nil && statusID != nil {
			return statusID
		}
	}

	return nil
}

// mapPaymentStatusID mapea el estado de pago (replicado para independencia de módulos)
func (uc *UseCaseUpdateOrder) mapPaymentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	if dto.PaymentStatusID != nil && *dto.PaymentStatusID > 0 {
		return dto.PaymentStatusID
	}

	for _, p := range dto.Payments {
		if p.Status == "completed" {
			paidID, err := uc.repo.GetPaymentStatusIDByCode(ctx, "paid")
			if err == nil && paidID != nil {
				return paidID
			}
			break
		}
	}

	if dto.IntegrationType == "shopify" && len(dto.PaymentDetails) > 0 {
		var paymentDetails map[string]interface{}
		if err := json.Unmarshal(dto.PaymentDetails, &paymentDetails); err != nil {
			return nil
		}

		financialStatus, ok := paymentDetails["financial_status"].(string)
		if !ok || financialStatus == "" {
			return nil
		}

		paymentStatusCode := statusmapper.MapShopifyFinancialStatusToPaymentStatus(financialStatus)
		mappedID, err := uc.repo.GetPaymentStatusIDByCode(ctx, paymentStatusCode)
		if err != nil || mappedID == nil {
			return nil
		}

		return mappedID
	}

	return nil
}

// mapFulfillmentStatusID mapea el estado de fulfillment (replicado para independencia de módulos)
func (uc *UseCaseUpdateOrder) mapFulfillmentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	if dto.FulfillmentStatusID != nil && *dto.FulfillmentStatusID > 0 {
		return dto.FulfillmentStatusID
	}

	if dto.IntegrationType == "shopify" && len(dto.FulfillmentDetails) > 0 {
		var fulfillmentDetails map[string]interface{}
		if err := json.Unmarshal(dto.FulfillmentDetails, &fulfillmentDetails); err != nil {
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatus, ok := fulfillmentDetails["fulfillment_status"].(string)
		if !ok || fulfillmentStatus == "" {
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatusCode := statusmapper.MapShopifyFulfillmentStatusToFulfillmentStatus(&fulfillmentStatus)
		return uc.getFulfillmentStatusIDByCode(ctx, fulfillmentStatusCode)
	}

	return nil
}

// getFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (uc *UseCaseUpdateOrder) getFulfillmentStatusIDByCode(ctx context.Context, code string) *uint {
	mappedID, err := uc.repo.GetFulfillmentStatusIDByCode(ctx, code)
	if err != nil || mappedID == nil {
		return nil
	}
	return mappedID
}

// syncIsPaidFromPaymentStatus sincroniza IsPaid basado en PaymentStatusID
func (uc *UseCaseUpdateOrder) syncIsPaidFromPaymentStatus(ctx context.Context, order *entities.ProbabilityOrder, paymentStatusID *uint) {
	if paymentStatusID == nil {
		return
	}

	paidStatusID, err := uc.repo.GetPaymentStatusIDByCode(ctx, "paid")
	if err != nil || paidStatusID == nil || *paidStatusID != *paymentStatusID {
		return
	}

	order.IsPaid = true
	if order.PaidAt == nil {
		now := time.Now()
		order.PaidAt = &now
	}
}
