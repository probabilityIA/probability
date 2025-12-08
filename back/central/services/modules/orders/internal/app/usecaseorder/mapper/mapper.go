package mapper

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// ToOrderResponse convierte un modelo Order a OrderResponse
func ToOrderResponse(order *domain.Order) *domain.OrderResponse {
	if order == nil {
		return nil
	}
	return &domain.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		// Identificadores de integración
		BusinessID:      order.BusinessID,
		IntegrationID:   order.IntegrationID,
		IntegrationType: order.IntegrationType,

		// Identificadores de la orden
		Platform:       order.Platform,
		ExternalID:     order.ExternalID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,

		// Información financiera
		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		CodTotal:     order.CodTotal,

		// Información del cliente
		CustomerID:    order.CustomerID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,

		// Dirección de envío (desnormalizado)
		ShippingStreet:     order.ShippingStreet,
		ShippingCity:       order.ShippingCity,
		ShippingState:      order.ShippingState,
		ShippingCountry:    order.ShippingCountry,
		ShippingPostalCode: order.ShippingPostalCode,
		ShippingLat:        order.ShippingLat,
		ShippingLng:        order.ShippingLng,

		// Información de pago
		PaymentMethodID: order.PaymentMethodID,
		IsPaid:          order.IsPaid,
		PaidAt:          order.PaidAt,

		// Información de envío/logística
		TrackingNumber: order.TrackingNumber,
		TrackingLink:   order.TrackingLink,
		GuideID:        order.GuideID,
		GuideLink:      order.GuideLink,
		DeliveryDate:   order.DeliveryDate,
		DeliveredAt:    order.DeliveredAt,

		// Información de fulfillment
		WarehouseID:   order.WarehouseID,
		WarehouseName: order.WarehouseName,
		DriverID:      order.DriverID,
		DriverName:    order.DriverName,
		IsLastMile:    order.IsLastMile,

		// Dimensiones y peso
		Weight: order.Weight,
		Height: order.Height,
		Width:  order.Width,
		Length: order.Length,
		Boxes:  order.Boxes,

		// Tipo y estado
		OrderTypeID:    order.OrderTypeID,
		OrderTypeName:  order.OrderTypeName,
		Status:         order.Status,
		OriginalStatus: order.OriginalStatus,

		// Información adicional
		Notes:    order.Notes,
		Coupon:   order.Coupon,
		Approved: order.Approved,
		UserID:   order.UserID,
		UserName: order.UserName,

		// Facturación
		Invoiceable:     order.Invoiceable,
		InvoiceURL:      order.InvoiceURL,
		InvoiceID:       order.InvoiceID,
		InvoiceProvider: order.InvoiceProvider,

		// Datos estructurados
		Items:              order.Items,
		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,

		// Timestamps
		OccurredAt: order.OccurredAt,
		ImportedAt: order.ImportedAt,
	}
}

// ToOrderSummary convierte un modelo Order a OrderSummary
func ToOrderSummary(order *domain.Order) domain.OrderSummary {
	paymentStatus := "unpaid"
	if order.IsPaid {
		paymentStatus = "paid"
	}

	var businessID uint
	if order.BusinessID != nil {
		businessID = *order.BusinessID
	}

	return domain.OrderSummary{
		ID:              order.ID,
		CreatedAt:       order.CreatedAt,
		BusinessID:      businessID,
		IntegrationID:   order.IntegrationID,
		IntegrationType: order.IntegrationType,
		Platform:        order.Platform,
		ExternalID:      order.ExternalID,
		OrderNumber:     order.OrderNumber,
		TotalAmount:     order.TotalAmount,
		Currency:        order.Currency,
		CustomerName:    order.CustomerName,
		CustomerEmail:   order.CustomerEmail,
		Status:          order.Status,
		PaymentStatus:   paymentStatus,
		ItemsCount:      len(order.Items),
	}
}
