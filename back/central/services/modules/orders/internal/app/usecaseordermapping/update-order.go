package usecaseordermapping

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"gorm.io/datatypes"
)

// UpdateOrder actualiza una orden existente con los datos del DTO
func (uc *UseCaseOrderMapping) UpdateOrder(ctx context.Context, existingOrder *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) (*domain.OrderResponse, error) {
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

	// 4. Publicar eventos relacionados con la actualización
	uc.publishUpdateEvents(ctx, existingOrder)

	// 5. Retornar la respuesta actualizada
	return uc.mapOrderToResponse(existingOrder), nil
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE ACTUALIZACIÓN DE CAMPOS
//
// ───────────────────────────────────────────

// updateOrderFields actualiza todos los campos de la orden y retorna si hubo cambios
func (uc *UseCaseOrderMapping) updateOrderFields(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
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

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE ESTADOS
//
// ───────────────────────────────────────────

// updateOrderStatuses actualiza los estados de la orden (OrderStatus, PaymentStatus, FulfillmentStatus)
func (uc *UseCaseOrderMapping) updateOrderStatuses(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
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

// updateOrderStatus actualiza el estado general de la orden
func (uc *UseCaseOrderMapping) updateOrderStatus(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	// Actualizar Status
	if dto.Status != "" && order.Status != dto.Status {
		order.Status = dto.Status
		changed = true
	}

	// Actualizar OriginalStatus y mapear StatusID
	if dto.OriginalStatus != "" && order.OriginalStatus != dto.OriginalStatus {
		order.OriginalStatus = dto.OriginalStatus
		changed = true

		// Buscar mapeo de estado cuando cambia el OriginalStatus
		mappedStatusID := uc.mapOrderStatusID(ctx, dto)
		if order.StatusID == nil || (mappedStatusID != nil && *order.StatusID != *mappedStatusID) || (mappedStatusID == nil && order.StatusID != nil) {
			order.StatusID = mappedStatusID
			changed = true
		}
	}

	return changed
}

// updatePaymentStatus actualiza el estado de pago
func (uc *UseCaseOrderMapping) updatePaymentStatus(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
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
func (uc *UseCaseOrderMapping) updateFulfillmentStatus(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
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

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS FINANCIEROS
//
// ───────────────────────────────────────────

// updateFinancialFields actualiza los campos financieros de la orden
func (uc *UseCaseOrderMapping) updateFinancialFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if dto.Subtotal > 0 && order.Subtotal != dto.Subtotal {
		order.Subtotal = dto.Subtotal
		changed = true
	}

	if dto.Tax >= 0 && order.Tax != dto.Tax {
		order.Tax = dto.Tax
		changed = true
	}

	if dto.Discount >= 0 && order.Discount != dto.Discount {
		order.Discount = dto.Discount
		changed = true
	}

	if dto.ShippingCost >= 0 && order.ShippingCost != dto.ShippingCost {
		order.ShippingCost = dto.ShippingCost
		changed = true
	}

	if dto.TotalAmount > 0 && order.TotalAmount != dto.TotalAmount {
		order.TotalAmount = dto.TotalAmount
		changed = true
	}

	if dto.Currency != "" && order.Currency != dto.Currency {
		order.Currency = dto.Currency
		changed = true
	}

	if dto.CodTotal != nil && (order.CodTotal == nil || *order.CodTotal != *dto.CodTotal) {
		order.CodTotal = dto.CodTotal
		changed = true
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS DEL CLIENTE
//
// ───────────────────────────────────────────

// updateCustomerFields actualiza la información del cliente
func (uc *UseCaseOrderMapping) updateCustomerFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if dto.CustomerName != "" && order.CustomerName != dto.CustomerName {
		order.CustomerName = dto.CustomerName
		changed = true
	}

	if dto.CustomerEmail != "" && order.CustomerEmail != dto.CustomerEmail {
		order.CustomerEmail = dto.CustomerEmail
		changed = true
	}

	if dto.CustomerPhone != "" && order.CustomerPhone != dto.CustomerPhone {
		order.CustomerPhone = dto.CustomerPhone
		changed = true
	}

	if dto.CustomerDNI != "" && order.CustomerDNI != dto.CustomerDNI {
		order.CustomerDNI = dto.CustomerDNI
		changed = true
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS DE ENVÍO
//
// ───────────────────────────────────────────

// updateShippingFields actualiza los campos relacionados con el envío
func (uc *UseCaseOrderMapping) updateShippingFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	// Actualizar información de tracking desde Shipments
	if uc.updateTrackingFields(order, dto) {
		changed = true
	}

	// Actualizar dirección de envío desde Addresses
	if uc.updateShippingAddress(order, dto) {
		changed = true
	}

	return changed
}

// updateTrackingFields actualiza los campos de tracking desde Shipments
func (uc *UseCaseOrderMapping) updateTrackingFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	if len(dto.Shipments) == 0 {
		return false
	}

	changed := false
	shipment := dto.Shipments[0]

	if shipment.TrackingNumber != nil && (order.TrackingNumber == nil || *order.TrackingNumber != *shipment.TrackingNumber) {
		order.TrackingNumber = shipment.TrackingNumber
		changed = true
	}

	if shipment.TrackingURL != nil && (order.TrackingLink == nil || *order.TrackingLink != *shipment.TrackingURL) {
		order.TrackingLink = shipment.TrackingURL
		changed = true
	}

	if shipment.GuideID != nil && (order.GuideID == nil || *order.GuideID != *shipment.GuideID) {
		order.GuideID = shipment.GuideID
		changed = true
	}

	if shipment.GuideURL != nil && (order.GuideLink == nil || *order.GuideLink != *shipment.GuideURL) {
		order.GuideLink = shipment.GuideURL
		changed = true
	}

	if shipment.DeliveredAt != nil && (order.DeliveredAt == nil || !order.DeliveredAt.Equal(*shipment.DeliveredAt)) {
		order.DeliveredAt = shipment.DeliveredAt
		changed = true
	}

	if shipment.ShippedAt != nil && (order.DeliveryDate == nil || !order.DeliveryDate.Equal(*shipment.ShippedAt)) {
		order.DeliveryDate = shipment.ShippedAt
		changed = true
	}

	return changed
}

// updateShippingAddress actualiza la dirección de envío desde Addresses
func (uc *UseCaseOrderMapping) updateShippingAddress(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	if len(dto.Addresses) == 0 {
		return false
	}

	changed := false

	for _, addr := range dto.Addresses {
		if addr.Type == "shipping" {
			if addr.Street != "" && order.ShippingStreet != addr.Street {
				order.ShippingStreet = addr.Street
				changed = true
			}

			if addr.Street2 != "" && order.Address2 != addr.Street2 {
				order.Address2 = addr.Street2
				changed = true
			}

			if addr.City != "" && order.ShippingCity != addr.City {
				order.ShippingCity = addr.City
				changed = true
			}

			if addr.State != "" && order.ShippingState != addr.State {
				order.ShippingState = addr.State
				changed = true
			}

			if addr.Country != "" && order.ShippingCountry != addr.Country {
				order.ShippingCountry = addr.Country
				changed = true
			}

			if addr.PostalCode != "" && order.ShippingPostalCode != addr.PostalCode {
				order.ShippingPostalCode = addr.PostalCode
				changed = true
			}

			if addr.Latitude != nil && (order.ShippingLat == nil || *order.ShippingLat != *addr.Latitude) {
				order.ShippingLat = addr.Latitude
				changed = true
			}

			if addr.Longitude != nil && (order.ShippingLng == nil || *order.ShippingLng != *addr.Longitude) {
				order.ShippingLng = addr.Longitude
				changed = true
			}

			break
		}
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS DE PAGO
//
// ───────────────────────────────────────────

// updatePaymentFields actualiza los campos relacionados con el pago
func (uc *UseCaseOrderMapping) updatePaymentFields(ctx context.Context, order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if len(dto.Payments) > 0 {
		payment := dto.Payments[0]

		if payment.PaymentMethodID > 0 && order.PaymentMethodID != payment.PaymentMethodID {
			order.PaymentMethodID = payment.PaymentMethodID
			changed = true
		}

		if payment.Status == "completed" && !order.IsPaid {
			order.IsPaid = true
			changed = true
		}

		if payment.PaidAt != nil && (order.PaidAt == nil || !order.PaidAt.Equal(*payment.PaidAt)) {
			order.PaidAt = payment.PaidAt
			changed = true
		}
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS DE FULFILLMENT
//
// ───────────────────────────────────────────

// updateFulfillmentFields actualiza los campos relacionados con fulfillment desde Shipments
func (uc *UseCaseOrderMapping) updateFulfillmentFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	if len(dto.Shipments) == 0 {
		return false
	}

	changed := false
	shipment := dto.Shipments[0]

	if shipment.WarehouseID != nil && (order.WarehouseID == nil || *order.WarehouseID != *shipment.WarehouseID) {
		order.WarehouseID = shipment.WarehouseID
		changed = true
	}

	if shipment.WarehouseName != "" && order.WarehouseName != shipment.WarehouseName {
		order.WarehouseName = shipment.WarehouseName
		changed = true
	}

	if shipment.DriverID != nil && (order.DriverID == nil || *order.DriverID != *shipment.DriverID) {
		order.DriverID = shipment.DriverID
		changed = true
	}

	if shipment.DriverName != "" && order.DriverName != shipment.DriverName {
		order.DriverName = shipment.DriverName
		changed = true
	}

	if order.IsLastMile != shipment.IsLastMile {
		order.IsLastMile = shipment.IsLastMile
		changed = true
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE CAMPOS ADICIONALES
//
// ───────────────────────────────────────────

// updateAdditionalFields actualiza campos adicionales de la orden
func (uc *UseCaseOrderMapping) updateAdditionalFields(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if dto.OrderTypeID != nil && (order.OrderTypeID == nil || *order.OrderTypeID != *dto.OrderTypeID) {
		order.OrderTypeID = dto.OrderTypeID
		changed = true
	}

	if dto.OrderTypeName != "" && order.OrderTypeName != dto.OrderTypeName {
		order.OrderTypeName = dto.OrderTypeName
		changed = true
	}

	if dto.Notes != nil && (order.Notes == nil || *order.Notes != *dto.Notes) {
		order.Notes = dto.Notes
		changed = true
	}

	if dto.Coupon != nil && (order.Coupon == nil || *order.Coupon != *dto.Coupon) {
		order.Coupon = dto.Coupon
		changed = true
	}

	if dto.Approved != nil && (order.Approved == nil || *order.Approved != *dto.Approved) {
		order.Approved = dto.Approved
		changed = true
	}

	if dto.UserID != nil && (order.UserID == nil || *order.UserID != *dto.UserID) {
		order.UserID = dto.UserID
		changed = true
	}

	if dto.UserName != "" && order.UserName != dto.UserName {
		order.UserName = dto.UserName
		changed = true
	}

	// Actualizar campos de facturación
	if order.Invoiceable != dto.Invoiceable {
		order.Invoiceable = dto.Invoiceable
		changed = true
	}

	if dto.InvoiceURL != nil && (order.InvoiceURL == nil || *order.InvoiceURL != *dto.InvoiceURL) {
		order.InvoiceURL = dto.InvoiceURL
		changed = true
	}

	if dto.InvoiceID != nil && (order.InvoiceID == nil || *order.InvoiceID != *dto.InvoiceID) {
		order.InvoiceID = dto.InvoiceID
		changed = true
	}

	if dto.InvoiceProvider != nil && (order.InvoiceProvider == nil || *order.InvoiceProvider != *dto.InvoiceProvider) {
		order.InvoiceProvider = dto.InvoiceProvider
		changed = true
	}

	if dto.OrderStatusURL != "" && order.OrderStatusURL != dto.OrderStatusURL {
		order.OrderStatusURL = dto.OrderStatusURL
		changed = true
	}

	return changed
}

// ───────────────────────────────────────────
//
//	ACTUALIZACIÓN DE DATOS ESTRUCTURADOS
//
// ───────────────────────────────────────────

// updateStructuredData actualiza los campos JSONB de la orden
func (uc *UseCaseOrderMapping) updateStructuredData(order *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	// Actualizar Items si están presentes
	if len(dto.Items) > 0 {
		if len(order.Items) == 0 || !equalJSON(order.Items, dto.Items) {
			order.Items = dto.Items
			changed = true
		}
	}

	// Actualizar Metadata si está presente
	if len(dto.Metadata) > 0 {
		if len(order.Metadata) == 0 || !equalJSON(order.Metadata, dto.Metadata) {
			order.Metadata = dto.Metadata
			changed = true
		}
	}

	// Actualizar FinancialDetails si está presente
	if len(dto.FinancialDetails) > 0 {
		if len(order.FinancialDetails) == 0 || !equalJSON(order.FinancialDetails, dto.FinancialDetails) {
			order.FinancialDetails = dto.FinancialDetails
			changed = true
		}
	}

	// Actualizar ShippingDetails si está presente
	if len(dto.ShippingDetails) > 0 {
		if len(order.ShippingDetails) == 0 || !equalJSON(order.ShippingDetails, dto.ShippingDetails) {
			order.ShippingDetails = dto.ShippingDetails
			changed = true
		}
	}

	// Actualizar PaymentDetails si está presente
	if len(dto.PaymentDetails) > 0 {
		if len(order.PaymentDetails) == 0 || !equalJSON(order.PaymentDetails, dto.PaymentDetails) {
			order.PaymentDetails = dto.PaymentDetails
			changed = true
		}
	}

	// Actualizar FulfillmentDetails si está presente
	if len(dto.FulfillmentDetails) > 0 {
		if len(order.FulfillmentDetails) == 0 || !equalJSON(order.FulfillmentDetails, dto.FulfillmentDetails) {
			order.FulfillmentDetails = dto.FulfillmentDetails
			changed = true
		}
	}

	return changed
}

// equalJSON compara dos valores JSONB
func equalJSON(a, b datatypes.JSON) bool {
	var aMap, bMap map[string]interface{}
	if err := json.Unmarshal(a, &aMap); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bMap); err != nil {
		return false
	}

	aBytes, _ := json.Marshal(aMap)
	bBytes, _ := json.Marshal(bMap)
	return string(aBytes) == string(bBytes)
}

// ───────────────────────────────────────────
//
//	FUNCIONES DE EVENTOS
//
// ───────────────────────────────────────────

// publishUpdateEvents publica los eventos relacionados con la actualización de la orden
func (uc *UseCaseOrderMapping) publishUpdateEvents(ctx context.Context, order *domain.ProbabilityOrder) {
	if uc.eventPublisher == nil {
		return
	}

	// Publicar evento de orden actualizada
	uc.publishOrderUpdatedEvent(ctx, order)

	// Recalcular score directamente
	uc.recalculateOrderScore(ctx, order)

	// Publicar evento para recalcular score (para otros consumidores)
	uc.publishScoreCalculationEventForUpdate(ctx, order)
}

// publishOrderUpdatedEvent publica el evento de orden actualizada
func (uc *UseCaseOrderMapping) publishOrderUpdatedEvent(ctx context.Context, order *domain.ProbabilityOrder) {
	eventData := domain.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
		CurrentStatus:  order.Status,
		CustomerEmail:  order.CustomerEmail,
		TotalAmount:    &order.TotalAmount,
		Currency:       order.Currency,
		Platform:       order.Platform,
	}

	event := domain.NewOrderEvent(domain.OrderEventTypeUpdated, order.ID, eventData)
	event.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		event.IntegrationID = &integrationID
	}

	go func() {
		if err := uc.eventPublisher.PublishOrderEvent(ctx, event); err != nil {
			uc.logger.Error(ctx).
				Err(err).
				Str("order_id", order.ID).
				Msg("Error al publicar evento de orden actualizada")
		}
	}()
}

// recalculateOrderScore recalcula el score de la orden
func (uc *UseCaseOrderMapping) recalculateOrderScore(ctx context.Context, order *domain.ProbabilityOrder) {
	go func() {
		fmt.Printf("[UpdateOrder] Recalculando score directamente para orden %s (actualizada)\n", order.ID)
		if err := uc.scoreUseCase.CalculateAndUpdateOrderScore(ctx, order.ID); err != nil {
			uc.logger.Error(ctx).
				Err(err).
				Str("order_id", order.ID).
				Msg("Error al recalcular score de la orden")
		} else {
			uc.logger.Info(ctx).
				Str("order_id", order.ID).
				Str("order_number", order.OrderNumber).
				Msg("✅ Score recalculado exitosamente para la orden actualizada")
		}
	}()
}

// publishScoreCalculationEventForUpdate publica el evento para recalcular score
func (uc *UseCaseOrderMapping) publishScoreCalculationEventForUpdate(ctx context.Context, order *domain.ProbabilityOrder) {
	scoreEventData := domain.OrderEventData{
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,
	}

	scoreEvent := domain.NewOrderEvent(domain.OrderEventTypeScoreCalculationRequested, order.ID, scoreEventData)
	scoreEvent.BusinessID = order.BusinessID
	if order.IntegrationID > 0 {
		integrationID := order.IntegrationID
		scoreEvent.IntegrationID = &integrationID
	}

	go func() {
		fmt.Printf("[UpdateOrder] Publicando evento order.score_calculation_requested para orden %s (actualizada)\n", order.ID)
		if err := uc.eventPublisher.PublishOrderEvent(ctx, scoreEvent); err != nil {
			uc.logger.Error(ctx).
				Err(err).
				Str("order_id", order.ID).
				Msg("Error al publicar evento de cálculo de score")
		} else {
			fmt.Printf("[UpdateOrder] Evento order.score_calculation_requested publicado exitosamente para orden %s\n", order.ID)
		}
	}()
}
