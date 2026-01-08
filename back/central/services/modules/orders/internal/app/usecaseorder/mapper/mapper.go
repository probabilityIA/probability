package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain" // Added import for service
	"gorm.io/datatypes"
)

// ToOrderResponse convierte un modelo Order a OrderResponse
func ToOrderResponse(order *domain.ProbabilityOrder) *domain.OrderResponse {
	if order == nil {
		return nil
	}

	// 1. Backward Compatibility: Populate Items (JSONB) from OrderItems (Relation) if Items (JSONB) is empty
	// This ensures that legacy orders or orders where JSONB wasn't populated still show products
	// 1. Prioritize OrderItems relation: Populate Items (JSONB) from OrderItems (Relation) if available
	// This ensures we serve the most up-to-date structured data from the order_items table
	items := order.Items
	if len(order.OrderItems) > 0 {
		if itemsJSON, err := json.Marshal(order.OrderItems); err == nil {
			items = datatypes.JSON(itemsJSON)
		}
	} else if len(items) == 0 || string(items) == "null" {
		// Fallback: If no OrderItems relation and no Items JSONB, ensure we return empty array, not null
		items = datatypes.JSON("[]")
	}

	// Checking imports first...
	return &domain.OrderResponse{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,

		// Identificadores de integración
		BusinessID:         order.BusinessID,
		IntegrationID:      order.IntegrationID,
		IntegrationType:    order.IntegrationType,
		IntegrationLogoURL: order.IntegrationLogoURL,

		// Identificadores de la orden
		Platform:       order.Platform,
		ExternalID:     order.ExternalID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,

		// Información financiera
		Subtotal:                order.Subtotal,
		Tax:                     order.Tax,
		Discount:                order.Discount,
		ShippingCost:            order.ShippingCost,
		TotalAmount:             order.TotalAmount,
		Currency:                order.Currency,
		CodTotal:                order.CodTotal,
		SubtotalPresentment:     order.SubtotalPresentment,
		TaxPresentment:          order.TaxPresentment,
		DiscountPresentment:     order.DiscountPresentment,
		ShippingCostPresentment: order.ShippingCostPresentment,
		TotalAmountPresentment:  order.TotalAmountPresentment,
		CurrencyPresentment:     order.CurrencyPresentment,

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
		TrackingNumber:      order.TrackingNumber,
		TrackingLink:        order.TrackingLink,
		GuideID:             order.GuideID,
		GuideLink:           order.GuideLink,
		DeliveryDate:        order.DeliveryDate,
		DeliveredAt:         order.DeliveredAt,
		DeliveryProbability: order.DeliveryProbability,

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
		OrderTypeID:         order.OrderTypeID,
		OrderTypeName:       order.OrderTypeName,
		Status:              order.Status,
		OriginalStatus:      order.OriginalStatus,
		StatusID:            order.StatusID,
		OrderStatus:         order.OrderStatus,
		PaymentStatusID:     order.PaymentStatusID,
		FulfillmentStatusID: order.FulfillmentStatusID,
		PaymentStatus:       order.PaymentStatus,
		FulfillmentStatus:   order.FulfillmentStatus,

		// Información adicional
		Notes:    order.Notes,
		Coupon:   order.Coupon,
		Approved: order.Approved,
		UserID:   order.UserID,
		UserName: order.UserName,

		// Novedades
		IsConfirmed: order.IsConfirmed,
		Novelty:     order.Novelty,

		// Facturación
		Invoiceable:     order.Invoiceable,
		InvoiceURL:      order.InvoiceURL,
		InvoiceID:       order.InvoiceID,
		InvoiceProvider: order.InvoiceProvider,
		OrderStatusURL:  order.OrderStatusURL,

		// Datos estructurados
		Items:              items,
		Metadata:           order.Metadata,
		FinancialDetails:   order.FinancialDetails,
		ShippingDetails:    order.ShippingDetails,
		PaymentDetails:     order.PaymentDetails,
		FulfillmentDetails: order.FulfillmentDetails,

		// Timestamps
		OccurredAt: order.OccurredAt,
		ImportedAt: order.ImportedAt,

		// Calculated Fields
		NegativeFactors: UnmarshalNegativeFactors(order.NegativeFactors),
	}
}

func UnmarshalNegativeFactors(jsonData datatypes.JSON) []string {
	if len(jsonData) == 0 || string(jsonData) == "null" {
		return []string{}
	}
	var factors []string
	_ = json.Unmarshal(jsonData, &factors)
	return factors
}

// ToOrderSummary convierte un modelo Order a OrderSummary
func ToOrderSummary(order *domain.ProbabilityOrder) domain.OrderSummary {
	var businessID uint
	if order.BusinessID != nil {
		businessID = *order.BusinessID
	}

	return domain.OrderSummary{
		ID:                     order.ID,
		CreatedAt:              order.CreatedAt,
		BusinessID:             businessID,
		IntegrationID:          order.IntegrationID,
		IntegrationType:        order.IntegrationType,
		IntegrationLogoURL:     order.IntegrationLogoURL,
		Platform:               order.Platform,
		ExternalID:             order.ExternalID,
		OrderNumber:            order.OrderNumber,
		TotalAmount:            order.TotalAmount,
		Currency:               order.Currency,
		TotalAmountPresentment: order.TotalAmountPresentment,
		CurrencyPresentment:    order.CurrencyPresentment,
		CustomerName:           order.CustomerName,
		CustomerEmail:          order.CustomerEmail,
		Status:                 order.Status,
		ItemsCount:             len(order.Items),
		DeliveryProbability:    order.DeliveryProbability,
		NegativeFactors:        UnmarshalNegativeFactors(order.NegativeFactors),
		OrderStatus:            order.OrderStatus,       // Información del estado de Probability
		PaymentStatus:          order.PaymentStatus,     // Información completa del estado de pago
		FulfillmentStatus:      order.FulfillmentStatus, // Información completa del estado de fulfillment
		IsConfirmed:            order.IsConfirmed,
		Novelty:                order.Novelty,
	}
}
