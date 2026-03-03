package usecasecreateorder

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/helpers/statusmapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// orderStatusMapping contiene los IDs de los estados mapeados
type orderStatusMapping struct {
	OrderStatusID       *uint
	PaymentStatusID     *uint
	FulfillmentStatusID *uint
}

// mapOrderStatuses mapea todos los estados de la orden (OrderStatus, PaymentStatus, FulfillmentStatus)
func (uc *UseCaseCreateOrder) mapOrderStatuses(ctx context.Context, dto *dtos.ProbabilityOrderDTO) orderStatusMapping {
	mapping := orderStatusMapping{}

	// Mapear OrderStatusID
	mapping.OrderStatusID = uc.mapOrderStatusID(ctx, dto)

	// Mapear PaymentStatusID
	mapping.PaymentStatusID = uc.mapPaymentStatusID(ctx, dto)

	// Mapear FulfillmentStatusID
	mapping.FulfillmentStatusID = uc.mapFulfillmentStatusID(ctx, dto)

	return mapping
}

// mapOrderStatusID mapea el estado general de la orden
// Prioridad 1: Intentar mapear usando order_status_mappings con Status
// Prioridad 2: Intentar mapear usando order_status_mappings con OriginalStatus
// Prioridad 3: Fallback directo por código en order_statuses (órdenes manuales o sin mapeo)
func (uc *UseCaseCreateOrder) mapOrderStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	// Intentar mapeo por integración si hay tipo de integración
	if dto.IntegrationType != "" {
		integrationTypeID := statusmapper.GetIntegrationTypeID(dto.IntegrationType)
		if integrationTypeID > 0 {
			// Prioridad 1: Intentar mapear usando Status
			if dto.Status != "" {
				mappedStatusID, err := uc.repo.GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx, integrationTypeID, dto.Status)
				if err == nil && mappedStatusID != nil {
					return mappedStatusID
				}
			}

			// Prioridad 2: Intentar con OriginalStatus (financial_status)
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

	// Prioridad 3: Fallback - buscar directamente por código en order_statuses
	// Útil para órdenes manuales o cuando no hay mapeo de integración configurado
	if dto.Status != "" {
		statusID, err := uc.repo.GetOrderStatusIDByCode(ctx, dto.Status)
		if err == nil && statusID != nil {
			return statusID
		}
	}

	return nil
}

// mapPaymentStatusID mapea el estado de pago desde el DTO o desde PaymentDetails si es Shopify
func (uc *UseCaseCreateOrder) mapPaymentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	// Si el DTO ya tiene PaymentStatusID, usarlo directamente
	if dto.PaymentStatusID != nil && *dto.PaymentStatusID > 0 {
		return dto.PaymentStatusID
	}

	// Para Shopify, extraer desde PaymentDetails
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

// mapFulfillmentStatusID mapea el estado de fulfillment desde el DTO o desde FulfillmentDetails si es Shopify
func (uc *UseCaseCreateOrder) mapFulfillmentStatusID(ctx context.Context, dto *dtos.ProbabilityOrderDTO) *uint {
	// Si el DTO ya tiene FulfillmentStatusID, usarlo directamente
	if dto.FulfillmentStatusID != nil && *dto.FulfillmentStatusID > 0 {
		return dto.FulfillmentStatusID
	}

	// Para Shopify, extraer desde FulfillmentDetails
	if dto.IntegrationType == "shopify" && len(dto.FulfillmentDetails) > 0 {
		var fulfillmentDetails map[string]interface{}
		if err := json.Unmarshal(dto.FulfillmentDetails, &fulfillmentDetails); err != nil {
			// Si no se puede parsear, usar "unfulfilled" por defecto
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatus, ok := fulfillmentDetails["fulfillment_status"].(string)
		if !ok || fulfillmentStatus == "" {
			// Si fulfillment_status es null o vacío, usar "unfulfilled"
			return uc.getFulfillmentStatusIDByCode(ctx, "unfulfilled")
		}

		fulfillmentStatusCode := statusmapper.MapShopifyFulfillmentStatusToFulfillmentStatus(&fulfillmentStatus)
		return uc.getFulfillmentStatusIDByCode(ctx, fulfillmentStatusCode)
	}

	return nil
}

// getFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (uc *UseCaseCreateOrder) getFulfillmentStatusIDByCode(ctx context.Context, code string) *uint {
	mappedID, err := uc.repo.GetFulfillmentStatusIDByCode(ctx, code)
	if err != nil || mappedID == nil {
		return nil
	}
	return mappedID
}

// syncIsPaidFromPaymentStatus sincroniza IsPaid basado en PaymentStatusID
func (uc *UseCaseCreateOrder) syncIsPaidFromPaymentStatus(ctx context.Context, order *entities.ProbabilityOrder, paymentStatusID *uint) {
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
