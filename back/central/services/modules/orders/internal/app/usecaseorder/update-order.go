package usecaseorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// UpdateOrder actualiza una orden existente
func (uc *UseCaseOrder) UpdateOrder(ctx context.Context, id string, req *domain.UpdateOrderRequest) (*domain.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	// Obtener la orden existente
	order, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	// Guardar el estado anterior para detectar cambios
	previousStatus := order.Status

	// Actualizar solo los campos proporcionados
	if req.Subtotal != nil {
		order.Subtotal = *req.Subtotal
	}
	if req.Tax != nil {
		order.Tax = *req.Tax
	}
	if req.Discount != nil {
		order.Discount = *req.Discount
	}
	if req.ShippingCost != nil {
		order.ShippingCost = *req.ShippingCost
	}
	if req.TotalAmount != nil {
		order.TotalAmount = *req.TotalAmount
	}
	if req.Currency != nil {
		order.Currency = *req.Currency
	}
	if req.CodTotal != nil {
		order.CodTotal = req.CodTotal
	}

	// Información del cliente
	if req.CustomerName != nil {
		order.CustomerName = *req.CustomerName
	}
	if req.CustomerEmail != nil {
		order.CustomerEmail = *req.CustomerEmail
	}
	if req.CustomerPhone != nil {
		order.CustomerPhone = *req.CustomerPhone
	}
	if req.CustomerDNI != nil {
		order.CustomerDNI = *req.CustomerDNI
	}

	// Dirección de envío
	if req.ShippingStreet != nil {
		order.ShippingStreet = *req.ShippingStreet
	}
	if req.ShippingCity != nil {
		order.ShippingCity = *req.ShippingCity
	}
	if req.ShippingState != nil {
		order.ShippingState = *req.ShippingState
	}
	if req.ShippingCountry != nil {
		order.ShippingCountry = *req.ShippingCountry
	}
	if req.ShippingPostalCode != nil {
		order.ShippingPostalCode = *req.ShippingPostalCode
	}
	if req.ShippingLat != nil {
		order.ShippingLat = req.ShippingLat
	}
	if req.ShippingLng != nil {
		order.ShippingLng = req.ShippingLng
	}

	// Información de pago
	if req.PaymentMethodID != nil {
		order.PaymentMethodID = *req.PaymentMethodID
	}
	if req.IsPaid != nil {
		order.IsPaid = *req.IsPaid
	}
	if req.PaidAt != nil {
		order.PaidAt = req.PaidAt
	}

	// Información de envío/logística
	if req.TrackingNumber != nil {
		order.TrackingNumber = req.TrackingNumber
	}
	if req.TrackingLink != nil {
		order.TrackingLink = req.TrackingLink
	}
	if req.GuideID != nil {
		order.GuideID = req.GuideID
	}
	if req.GuideLink != nil {
		order.GuideLink = req.GuideLink
	}
	if req.DeliveryDate != nil {
		order.DeliveryDate = req.DeliveryDate
	}
	if req.DeliveredAt != nil {
		order.DeliveredAt = req.DeliveredAt
	}

	// Información de fulfillment
	if req.WarehouseID != nil {
		order.WarehouseID = req.WarehouseID
	}
	if req.WarehouseName != nil {
		order.WarehouseName = *req.WarehouseName
	}
	if req.DriverID != nil {
		order.DriverID = req.DriverID
	}
	if req.DriverName != nil {
		order.DriverName = *req.DriverName
	}
	if req.IsLastMile != nil {
		order.IsLastMile = *req.IsLastMile
	}

	// Dimensiones y peso
	if req.Weight != nil {
		order.Weight = req.Weight
	}
	if req.Height != nil {
		order.Height = req.Height
	}
	if req.Width != nil {
		order.Width = req.Width
	}
	if req.Length != nil {
		order.Length = req.Length
	}
	if req.Boxes != nil {
		order.Boxes = req.Boxes
	}

	// Tipo y estado
	if req.OrderTypeID != nil {
		order.OrderTypeID = req.OrderTypeID
	}
	if req.OrderTypeName != nil {
		order.OrderTypeName = *req.OrderTypeName
	}
	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.OriginalStatus != nil {
		order.OriginalStatus = *req.OriginalStatus
	}

	// Información adicional
	if req.Notes != nil {
		order.Notes = req.Notes
	}
	if req.Coupon != nil {
		order.Coupon = req.Coupon
	}
	if req.Approved != nil {
		order.Approved = req.Approved
	}
	if req.UserID != nil {
		order.UserID = req.UserID
	}
	if req.UserName != nil {
		order.UserName = *req.UserName
	}

	// Facturación
	if req.Invoiceable != nil {
		order.Invoiceable = *req.Invoiceable
	}
	if req.InvoiceURL != nil {
		order.InvoiceURL = req.InvoiceURL
	}
	if req.InvoiceID != nil {
		order.InvoiceID = req.InvoiceID
	}
	if req.InvoiceProvider != nil {
		order.InvoiceProvider = req.InvoiceProvider
	}

	// Datos estructurados
	if req.Items != nil {
		order.Items = req.Items
	}
	if req.Metadata != nil {
		order.Metadata = req.Metadata
	}
	if req.FinancialDetails != nil {
		order.FinancialDetails = req.FinancialDetails
	}
	if req.ShippingDetails != nil {
		order.ShippingDetails = req.ShippingDetails
	}
	if req.PaymentDetails != nil {
		order.PaymentDetails = req.PaymentDetails
	}
	if req.FulfillmentDetails != nil {
		order.FulfillmentDetails = req.FulfillmentDetails
	}

	// Guardar cambios
	if err := uc.repo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	// Publicar eventos si hay publicador disponible
	if uc.eventPublisher != nil {
		// Publicar evento de actualización
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
		// Publicar de forma asíncrona
		go func() {
			uc.eventPublisher.PublishOrderEvent(ctx, event)
		}()

		// Si cambió el estado, publicar evento de cambio de estado
		if previousStatus != order.Status {
			statusEventData := domain.OrderEventData{
				OrderNumber:    order.OrderNumber,
				InternalNumber: order.InternalNumber,
				ExternalID:     order.ExternalID,
				PreviousStatus: previousStatus,
				CurrentStatus:  order.Status,
				CustomerEmail:  order.CustomerEmail,
				TotalAmount:    &order.TotalAmount,
				Currency:       order.Currency,
				Platform:       order.Platform,
			}
			statusEvent := domain.NewOrderEvent(domain.OrderEventTypeStatusChanged, order.ID, statusEventData)
			statusEvent.BusinessID = order.BusinessID
			if order.IntegrationID > 0 {
				integrationID := order.IntegrationID
				statusEvent.IntegrationID = &integrationID
			}
			// Publicar de forma asíncrona
			go func() {
				uc.eventPublisher.PublishOrderEvent(ctx, statusEvent)
			}()
		}
	}

	return mapper.ToOrderResponse(order), nil
}
