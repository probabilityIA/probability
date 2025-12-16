package usecaseordermapping

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

func (uc *UseCaseOrderMapping) UpdateOrder(ctx context.Context, existingOrder *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) (*domain.OrderResponse, error) {
	hasChanges := false

	if validateAndUpdateStatus(existingOrder, dto) {
		hasChanges = true
	}

	if validateAndUpdateShipping(existingOrder, dto) {
		hasChanges = true
	}

	if validateAndUpdatePayment(existingOrder, dto) {
		hasChanges = true
	}

	if !hasChanges {
		return uc.mapOrderToResponse(existingOrder), nil
	}

	if err := uc.repo.UpdateOrder(ctx, existingOrder); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	if uc.eventPublisher != nil {
		// Publicar evento de orden actualizada
		eventData := domain.OrderEventData{
			OrderNumber:    existingOrder.OrderNumber,
			InternalNumber: existingOrder.InternalNumber,
			ExternalID:     existingOrder.ExternalID,
			CurrentStatus:  existingOrder.Status,
			CustomerEmail:  existingOrder.CustomerEmail,
			TotalAmount:    &existingOrder.TotalAmount,
			Currency:       existingOrder.Currency,
			Platform:       existingOrder.Platform,
		}
		event := domain.NewOrderEvent(domain.OrderEventTypeUpdated, existingOrder.ID, eventData)
		event.BusinessID = existingOrder.BusinessID
		if existingOrder.IntegrationID > 0 {
			integrationID := existingOrder.IntegrationID
			event.IntegrationID = &integrationID
		}
		go func() {
			if err := uc.eventPublisher.PublishOrderEvent(ctx, event); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", existingOrder.ID).
					Msg("Error al publicar evento de orden actualizada")
			}
		}()

		// Recalcular score directamente cuando se actualiza (porque los datos pueden haber cambiado)
		go func() {
			fmt.Printf("[UpdateOrder] Recalculando score directamente para orden %s (actualizada)\n", existingOrder.ID)
			if err := uc.scoreUseCase.CalculateAndUpdateOrderScore(ctx, existingOrder.ID); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", existingOrder.ID).
					Msg("Error al recalcular score de la orden")
			} else {
				uc.logger.Info(ctx).
					Str("order_id", existingOrder.ID).
					Str("order_number", existingOrder.OrderNumber).
					Msg("✅ Score recalculado exitosamente para la orden actualizada")
			}
		}()

		// Publicar evento para recalcular score (mantener para otros consumidores)
		scoreEventData := domain.OrderEventData{
			OrderNumber:    existingOrder.OrderNumber,
			InternalNumber: existingOrder.InternalNumber,
			ExternalID:     existingOrder.ExternalID,
		}
		scoreEvent := domain.NewOrderEvent(domain.OrderEventTypeScoreCalculationRequested, existingOrder.ID, scoreEventData)
		scoreEvent.BusinessID = existingOrder.BusinessID
		if existingOrder.IntegrationID > 0 {
			integrationID := existingOrder.IntegrationID
			scoreEvent.IntegrationID = &integrationID
		}
		go func() {
			fmt.Printf("[UpdateOrder] Publicando evento order.score_calculation_requested para orden %s (actualizada)\n", existingOrder.ID)
			if err := uc.eventPublisher.PublishOrderEvent(ctx, scoreEvent); err != nil {
				uc.logger.Error(ctx).
					Err(err).
					Str("order_id", existingOrder.ID).
					Msg("Error al publicar evento de cálculo de score")
			} else {
				fmt.Printf("[UpdateOrder] Evento order.score_calculation_requested publicado exitosamente para orden %s\n", existingOrder.ID)
			}
		}()
	}

	return uc.mapOrderToResponse(existingOrder), nil
}

func validateAndUpdateStatus(existingOrder *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if dto.Status != "" && existingOrder.Status != dto.Status {
		existingOrder.Status = dto.Status
		changed = true
	}

	if dto.OriginalStatus != "" && existingOrder.OriginalStatus != dto.OriginalStatus {
		existingOrder.OriginalStatus = dto.OriginalStatus
		changed = true
	}

	return changed
}

func validateAndUpdateShipping(existingOrder *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if len(dto.Shipments) > 0 {
		shipment := dto.Shipments[0]
		if shipment.TrackingNumber != nil && (existingOrder.TrackingNumber == nil || *existingOrder.TrackingNumber != *shipment.TrackingNumber) {
			existingOrder.TrackingNumber = shipment.TrackingNumber
			changed = true
		}
		if shipment.TrackingURL != nil && (existingOrder.TrackingLink == nil || *existingOrder.TrackingLink != *shipment.TrackingURL) {
			existingOrder.TrackingLink = shipment.TrackingURL
			changed = true
		}
		if shipment.GuideID != nil && (existingOrder.GuideID == nil || *existingOrder.GuideID != *shipment.GuideID) {
			existingOrder.GuideID = shipment.GuideID
			changed = true
		}
		if shipment.GuideURL != nil && (existingOrder.GuideLink == nil || *existingOrder.GuideLink != *shipment.GuideURL) {
			existingOrder.GuideLink = shipment.GuideURL
			changed = true
		}
		if shipment.DeliveredAt != nil && (existingOrder.DeliveredAt == nil || !existingOrder.DeliveredAt.Equal(*shipment.DeliveredAt)) {
			existingOrder.DeliveredAt = shipment.DeliveredAt
			changed = true
		}
		if shipment.ShippedAt != nil {
			if existingOrder.DeliveryDate == nil || !existingOrder.DeliveryDate.Equal(*shipment.ShippedAt) {
				existingOrder.DeliveryDate = shipment.ShippedAt
				changed = true
			}
		}
	}

	if len(dto.Addresses) > 0 {
		for _, addr := range dto.Addresses {
			if addr.Type == "shipping" {
				if addr.Street != "" && existingOrder.ShippingStreet != addr.Street {
					existingOrder.ShippingStreet = addr.Street
					changed = true
				}
				if addr.Street2 != "" {
					existingOrder.ShippingStreet2 = addr.Street2
					changed = true
				}
				if addr.City != "" && existingOrder.ShippingCity != addr.City {
					existingOrder.ShippingCity = addr.City
					changed = true
				}
				if addr.State != "" && existingOrder.ShippingState != addr.State {
					existingOrder.ShippingState = addr.State
					changed = true
				}
				if addr.Country != "" && existingOrder.ShippingCountry != addr.Country {
					existingOrder.ShippingCountry = addr.Country
					changed = true
				}
				if addr.PostalCode != "" && existingOrder.ShippingPostalCode != addr.PostalCode {
					existingOrder.ShippingPostalCode = addr.PostalCode
					changed = true
				}
				if addr.Latitude != nil && (existingOrder.ShippingLat == nil || *existingOrder.ShippingLat != *addr.Latitude) {
					existingOrder.ShippingLat = addr.Latitude
					changed = true
				}
				if addr.Longitude != nil && (existingOrder.ShippingLng == nil || *existingOrder.ShippingLng != *addr.Longitude) {
					existingOrder.ShippingLng = addr.Longitude
					changed = true
				}
				break
			}
		}
	}

	return changed
}

func validateAndUpdatePayment(existingOrder *domain.ProbabilityOrder, dto *domain.ProbabilityOrderDTO) bool {
	changed := false

	if len(dto.Payments) > 0 {
		payment := dto.Payments[0]
		if payment.PaymentMethodID > 0 && existingOrder.PaymentMethodID != payment.PaymentMethodID {
			existingOrder.PaymentMethodID = payment.PaymentMethodID
			changed = true
		}
		if payment.Status == "completed" && !existingOrder.IsPaid {
			existingOrder.IsPaid = true
			changed = true
		}
		if payment.PaidAt != nil && (existingOrder.PaidAt == nil || !existingOrder.PaidAt.Equal(*payment.PaidAt)) {
			existingOrder.PaidAt = payment.PaidAt
			changed = true
		}
	}

	return changed
}

func (uc *UseCaseOrderMapping) mapOrderToResponse(order *domain.ProbabilityOrder) *domain.OrderResponse {
	return &domain.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		BusinessID:         order.BusinessID,
		IntegrationID:      order.IntegrationID,
		IntegrationType:    order.IntegrationType,
		IntegrationLogoURL: order.IntegrationLogoURL,

		Platform:       order.Platform,
		ExternalID:     order.ExternalID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,

		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		CodTotal:     order.CodTotal,

		CustomerID:    order.CustomerID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,

		ShippingStreet:     order.ShippingStreet,
		ShippingCity:       order.ShippingCity,
		ShippingState:      order.ShippingState,
		ShippingCountry:    order.ShippingCountry,
		ShippingPostalCode: order.ShippingPostalCode,
		ShippingLat:        order.ShippingLat,
		ShippingLng:        order.ShippingLng,

		PaymentMethodID: order.PaymentMethodID,
		IsPaid:          order.IsPaid,
		PaidAt:          order.PaidAt,

		TrackingNumber:      order.TrackingNumber,
		TrackingLink:        order.TrackingLink,
		GuideID:             order.GuideID,
		GuideLink:           order.GuideLink,
		DeliveryDate:        order.DeliveryDate,
		DeliveredAt:         order.DeliveredAt,
		DeliveryProbability: order.DeliveryProbability,

		WarehouseID:   order.WarehouseID,
		WarehouseName: order.WarehouseName,
		DriverID:      order.DriverID,
		DriverName:    order.DriverName,
		IsLastMile:    order.IsLastMile,

		Weight: order.Weight,
		Height: order.Height,
		Width:  order.Width,
		Length: order.Length,
		Boxes:  order.Boxes,

		OrderTypeID:    order.OrderTypeID,
		OrderTypeName:  order.OrderTypeName,
		Status:         order.Status,
		OriginalStatus: order.OriginalStatus,

		Notes:    order.Notes,
		Coupon:   order.Coupon,
		Approved: order.Approved,
		UserID:   order.UserID,
		UserName: order.UserName,

		Invoiceable:     order.Invoiceable,
		InvoiceURL:      order.InvoiceURL,
		InvoiceID:       order.InvoiceID,
		InvoiceProvider: order.InvoiceProvider,
		OrderStatusURL:  order.OrderStatusURL,

		Items:              order.Items,
		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,

		OccurredAt: order.OccurredAt,
		ImportedAt: order.ImportedAt,
	}
}
