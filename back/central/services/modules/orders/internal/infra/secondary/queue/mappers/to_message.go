package mappers

import (
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/response"
)

func OrderToSnapshot(order *entities.ProbabilityOrder) *response.OrderSnapshot {
	return &response.OrderSnapshot{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,

		TotalAmount:     order.TotalAmount,
		Currency:        order.Currency,
		PaymentMethodID: order.PaymentMethodID,
		PaymentStatusID: order.PaymentStatusID,

		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,

		CustomerID:    order.CustomerID,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,

		Platform:      order.Platform,
		IntegrationID: order.IntegrationID,
		BusinessName:  order.BusinessName,

		OrderStatusID:       order.StatusID,
		FulfillmentStatusID: order.FulfillmentStatusID,

		WarehouseID: order.WarehouseID,

		Items: OrderItemsToSnapshot(order.OrderItems),

		ShippingStreet:      order.ShippingStreet,
		ShippingCity:        order.ShippingCity,
		ShippingState:       order.ShippingState,
		ShippingCountry:     order.ShippingCountry,
		ShippingPostalCode:  order.ShippingPostalCode,
		ShippingLat:         order.ShippingLat,
		ShippingLng:         order.ShippingLng,
		ItemsSummary:        BuildItemsSummary(order.OrderItems),
		ShippingAddress:     BuildShippingAddress(order),
		PaymentMethodName:   order.PaymentMethodName,
		TrackingNumber:      derefString(order.TrackingNumber),
		Carrier:             extractCarrier(order),
		IsPaid:              order.IsPaid,
		DeliveryProbability: derefFloat64(order.DeliveryProbability),
		Status:              order.Status,

		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func OrderItemsToSnapshot(items []entities.ProbabilityOrderItem) []response.OrderItemSnapshot {
	if len(items) == 0 {
		return []response.OrderItemSnapshot{}
	}

	snapshots := make([]response.OrderItemSnapshot, 0, len(items))
	for _, item := range items {
		name := item.ProductName
		if name == "" {
			name = item.ProductTitle
		}
		if name == "" {
			name = item.ProductSKU
		}

		snapshot := response.OrderItemSnapshot{
			ProductID:       item.ProductID,
			SKU:             item.ProductSKU,
			VariantID:       item.VariantID,
			Name:            name,
			Title:           item.ProductTitle,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			TotalPrice:      item.TotalPrice,
			Tax:             item.Tax,
			TaxRate:         item.TaxRate,
			Discount:        item.Discount,
			DiscountPercent: item.DiscountPercent,
			ImageURL:        item.ImageURL,
			ProductURL:      item.ProductURL,
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}

func BuildItemsSummary(items []entities.ProbabilityOrderItem) string {
	if len(items) == 0 {
		return "Sin items"
	}

	var parts []string
	for _, item := range items {
		productName := item.ProductName
		if productName == "" {
			productName = item.ProductTitle
		}
		if productName == "" {
			productName = item.ProductSKU
		}
		if productName == "" {
			productName = "Producto sin nombre"
		}

		parts = append(parts, fmt.Sprintf("%dx %s", item.Quantity, productName))
	}

	return strings.Join(parts, ", ")
}

func BuildShippingAddress(order *entities.ProbabilityOrder) string {
	var parts []string

	if order.ShippingStreet != "" {
		parts = append(parts, order.ShippingStreet)
	}
	if order.ShippingCity != "" {
		parts = append(parts, order.ShippingCity)
	}
	if order.ShippingState != "" {
		parts = append(parts, order.ShippingState)
	}
	if order.ShippingCountry != "" {
		parts = append(parts, order.ShippingCountry)
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}

func EventToMessage(event *entities.OrderEvent) *response.OrderEventMessage {
	return &response.OrderEventMessage{
		EventID:       event.ID,
		EventType:     string(event.Type),
		OrderID:       event.OrderID,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		Timestamp:     event.Timestamp,
		Changes: map[string]any{
			"previous_status": event.Data.PreviousStatus,
			"current_status":  event.Data.CurrentStatus,
			"platform":        event.Data.Platform,
		},
		Metadata: event.Metadata,
	}
}

func GenerateEventID() string {
	return time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
}

func derefFloat64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func extractCarrier(order *entities.ProbabilityOrder) string {
	if len(order.Shipments) > 0 && order.Shipments[0].Carrier != nil {
		return *order.Shipments[0].Carrier
	}
	return ""
}
