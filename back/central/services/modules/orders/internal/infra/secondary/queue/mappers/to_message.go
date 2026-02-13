package mappers

import (
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/response"
)

// OrderToSnapshot convierte una ProbabilityOrder a OrderSnapshot COMPLETO
// Incluye toda la información necesaria para que consumidores externos puedan
// decidir si actuar sin necesidad de consultar la base de datos
func OrderToSnapshot(order *entities.ProbabilityOrder) *response.OrderSnapshot {
	return &response.OrderSnapshot{
		// Identificadores
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		InternalNumber: order.InternalNumber,
		ExternalID:     order.ExternalID,

		// Información financiera
		TotalAmount:     order.TotalAmount,
		Currency:        order.Currency,
		PaymentMethodID: order.PaymentMethodID,
		PaymentStatusID: order.PaymentStatusID,

		// Información financiera detallada (para facturación)
		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,

		// Información del cliente
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,

		// Información de origen
		Platform:      order.Platform,
		IntegrationID: order.IntegrationID,

		// Estados
		OrderStatusID:       order.StatusID,
		FulfillmentStatusID: order.FulfillmentStatusID,

		// Items detallados (para facturación e inventario)
		Items: OrderItemsToSnapshot(order.OrderItems),

		// Información adicional para mensajes
		ItemsSummary:    BuildItemsSummary(order.OrderItems),
		ShippingAddress: BuildShippingAddress(order),

		// Timestamps
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

// OrderItemsToSnapshot convierte un slice de ProbabilityOrderItem a OrderItemSnapshot
// Mapea toda la información de los items necesaria para facturación e inventario
func OrderItemsToSnapshot(items []entities.ProbabilityOrderItem) []response.OrderItemSnapshot {
	if len(items) == 0 {
		return []response.OrderItemSnapshot{}
	}

	snapshots := make([]response.OrderItemSnapshot, 0, len(items))
	for _, item := range items {
		// Elegir el mejor nombre disponible
		name := item.ProductName
		if name == "" {
			name = item.ProductTitle
		}
		if name == "" {
			name = item.ProductSKU
		}

		// Crear snapshot del item
		snapshot := response.OrderItemSnapshot{
			ProductID: item.ProductID,
			SKU:       item.ProductSKU,
			VariantID: item.VariantID,
			Name:      name,
			Title:     item.ProductTitle,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			TotalPrice: item.TotalPrice,
			Tax:       item.Tax,
			TaxRate:   item.TaxRate,
			Discount:  item.Discount,
			ImageURL:  item.ImageURL,
			ProductURL: item.ProductURL,
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}

// BuildItemsSummary construye un resumen legible de los items de una orden
// Formato: "2x Producto A, 1x Producto B, 3x Producto C"
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

// BuildShippingAddress construye una dirección legible de envío
// Formato: "Calle 123 #45-67, Bogotá, Cundinamarca, Colombia"
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

// EventToMessage convierte un OrderEvent a OrderEventMessage
func EventToMessage(event *entities.OrderEvent) *response.OrderEventMessage {
	return &response.OrderEventMessage{
		EventID:       event.ID,
		EventType:     string(event.Type),
		OrderID:       event.OrderID,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		Timestamp:     event.Timestamp,
		Changes: map[string]interface{}{
			"previous_status": event.Data.PreviousStatus,
			"current_status":  event.Data.CurrentStatus,
			"platform":        event.Data.Platform,
		},
		Metadata: event.Metadata,
	}
}

// GenerateEventID genera un ID único para el evento
func GenerateEventID() string {
	return time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
}
