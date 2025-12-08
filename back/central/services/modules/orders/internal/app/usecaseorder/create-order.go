package usecaseorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// CreateOrder crea una nueva orden
func (uc *UseCaseOrder) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	// Validar que no exista una orden con el mismo external_id para la misma integración
	exists, err := uc.repo.OrderExists(ctx, req.ExternalID, req.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error checking if order exists: %w", err)
	}
	if exists {
		return nil, errors.New("order with this external_id already exists for this integration")
	}

	// Crear el modelo de orden
	order := &domain.Order{
		// Identificadores de integración
		BusinessID:      req.BusinessID,
		IntegrationID:   req.IntegrationID,
		IntegrationType: req.IntegrationType,

		// Identificadores de la orden
		Platform:       req.Platform,
		ExternalID:     req.ExternalID,
		OrderNumber:    req.OrderNumber,
		InternalNumber: req.InternalNumber,

		// Información financiera
		Subtotal:     req.Subtotal,
		Tax:          req.Tax,
		Discount:     req.Discount,
		ShippingCost: req.ShippingCost,
		TotalAmount:  req.TotalAmount,
		Currency:     req.Currency,
		CodTotal:     req.CodTotal,

		// Información del cliente
		CustomerID:    req.CustomerID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		CustomerDNI:   req.CustomerDNI,

		// Dirección de envío
		ShippingStreet:     req.ShippingStreet,
		ShippingCity:       req.ShippingCity,
		ShippingState:      req.ShippingState,
		ShippingCountry:    req.ShippingCountry,
		ShippingPostalCode: req.ShippingPostalCode,
		ShippingLat:        req.ShippingLat,
		ShippingLng:        req.ShippingLng,

		// Información de pago
		PaymentMethodID: req.PaymentMethodID,
		IsPaid:          req.IsPaid,
		PaidAt:          req.PaidAt,

		// Información de envío/logística
		TrackingNumber: req.TrackingNumber,
		TrackingLink:   req.TrackingLink,
		GuideID:        req.GuideID,
		GuideLink:      req.GuideLink,
		DeliveryDate:   req.DeliveryDate,
		DeliveredAt:    req.DeliveredAt,

		// Información de fulfillment
		WarehouseID:   req.WarehouseID,
		WarehouseName: req.WarehouseName,
		DriverID:      req.DriverID,
		DriverName:    req.DriverName,
		IsLastMile:    req.IsLastMile,

		// Dimensiones y peso
		Weight: req.Weight,
		Height: req.Height,
		Width:  req.Width,
		Length: req.Length,
		Boxes:  req.Boxes,

		// Tipo y estado
		OrderTypeID:    req.OrderTypeID,
		OrderTypeName:  req.OrderTypeName,
		Status:         req.Status,
		OriginalStatus: req.OriginalStatus,

		// Información adicional
		Notes:    req.Notes,
		Coupon:   req.Coupon,
		Approved: req.Approved,
		UserID:   req.UserID,
		UserName: req.UserName,

		// Facturación
		Invoiceable:     req.Invoiceable,
		InvoiceURL:      req.InvoiceURL,
		InvoiceID:       req.InvoiceID,
		InvoiceProvider: req.InvoiceProvider,

		// Datos estructurados
		Items:              req.Items,
		Metadata:           req.Metadata,
		FinancialDetails:   req.FinancialDetails,
		ShippingDetails:    req.ShippingDetails,
		PaymentDetails:     req.PaymentDetails,
		FulfillmentDetails: req.FulfillmentDetails,

		// Timestamps
		OccurredAt: req.OccurredAt,
		ImportedAt: req.ImportedAt,
	}

	// Guardar en la base de datos
	if err := uc.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Publicar evento de orden creada
	if uc.eventPublisher != nil {
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
		event := domain.NewOrderEvent(domain.OrderEventTypeCreated, order.ID, eventData)
		event.BusinessID = order.BusinessID
		if order.IntegrationID > 0 {
			integrationID := order.IntegrationID
			event.IntegrationID = &integrationID
		}
		// Publicar de forma asíncrona (no bloquear si falla)
		go func() {
			if err := uc.eventPublisher.PublishOrderEvent(ctx, event); err != nil {
				// Log error pero no fallar la creación de la orden
			}
		}()
	}

	// Retornar la respuesta
	return mapper.ToOrderResponse(order), nil
}
