package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/response"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// OrderRabbitPublisher implementa el publicador de eventos de órdenes
type OrderRabbitPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewOrderRabbitPublisher crea un nuevo publisher de RabbitMQ
func NewOrderRabbitPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) ports.IOrderRabbitPublisher {
	return &OrderRabbitPublisher{
		rabbit: rabbit,
		log:    logger,
	}
}

// ───────────────────────────────────────────
//
//	MÉTODOS DE PUBLICACIÓN DE EVENTOS
//
// ───────────────────────────────────────────

// PublishOrderCreated publica un evento de orden creada
func (p *OrderRabbitPublisher) PublishOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) error {
	message := &response.OrderEventMessage{
		EventID:       mappers.GenerateEventID(),
		EventType:     "order.created",
		OrderID:       order.ID,
		BusinessID:    order.BusinessID,
		IntegrationID: &order.IntegrationID,
		Timestamp:     time.Now(),
		Order:         mappers.OrderToSnapshot(order),
	}
	return p.publishToQueue(ctx, "orders.events.created", message)
}

// PublishOrderUpdated publica un evento de orden actualizada
func (p *OrderRabbitPublisher) PublishOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) error {
	message := &response.OrderEventMessage{
		EventID:       mappers.GenerateEventID(),
		EventType:     "order.updated",
		OrderID:       order.ID,
		BusinessID:    order.BusinessID,
		IntegrationID: &order.IntegrationID,
		Timestamp:     time.Now(),
		Order:         mappers.OrderToSnapshot(order),
	}
	return p.publishToQueue(ctx, "orders.events.updated", message)
}

// PublishOrderCancelled publica un evento de orden cancelada
func (p *OrderRabbitPublisher) PublishOrderCancelled(ctx context.Context, order *entities.ProbabilityOrder) error {
	message := &response.OrderEventMessage{
		EventID:       mappers.GenerateEventID(),
		EventType:     "order.cancelled",
		OrderID:       order.ID,
		BusinessID:    order.BusinessID,
		IntegrationID: &order.IntegrationID,
		Timestamp:     time.Now(),
		Order:         mappers.OrderToSnapshot(order),
	}
	return p.publishToQueue(ctx, "orders.events.cancelled", message)
}

// PublishOrderStatusChanged publica un evento de cambio de estado
func (p *OrderRabbitPublisher) PublishOrderStatusChanged(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error {
	message := &response.OrderEventMessage{
		EventID:       mappers.GenerateEventID(),
		EventType:     "order.status_changed",
		OrderID:       order.ID,
		BusinessID:    order.BusinessID,
		IntegrationID: &order.IntegrationID,
		Timestamp:     time.Now(),
		Changes: map[string]interface{}{
			"previous_status": previousStatus,
			"current_status":  currentStatus,
		},
		Order: mappers.OrderToSnapshot(order),
	}
	return p.publishToQueue(ctx, "orders.events.status_changed", message)
}

// PublishOrderEvent publica un evento genérico de orden con snapshot completo
func (p *OrderRabbitPublisher) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	queue := p.getQueueForEventType(event.Type)

	// Construir mensaje completo con OrderSnapshot
	message := &response.OrderEventMessage{
		EventID:       event.ID,
		EventType:     string(event.Type),
		OrderID:       event.OrderID,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		Timestamp:     event.Timestamp,
		Order:         mappers.OrderToSnapshot(order), // ✅ Snapshot completo
		Changes: map[string]interface{}{
			"previous_status": event.Data.PreviousStatus,
			"current_status":  event.Data.CurrentStatus,
			"platform":        event.Data.Platform,
		},
		Metadata: event.Metadata,
	}

	return p.publishToQueue(ctx, queue, message)
}

// PublishConfirmationRequested publica un evento cuando una orden requiere confirmación
func (p *OrderRabbitPublisher) PublishConfirmationRequested(ctx context.Context, order *entities.ProbabilityOrder) error {
	// Construir resumen de items
	itemsSummary := buildItemsSummary(order)

	// Construir dirección de envío
	shippingAddress := buildShippingAddress(order)

	// Construir evento
	event := map[string]interface{}{
		"event_type":        "order.confirmation_requested",
		"order_id":          order.ID,
		"order_number":      order.OrderNumber,
		"business_id":       order.BusinessID,
		"customer_name":     order.CustomerName,
		"customer_phone":    order.CustomerPhone,
		"customer_email":    order.CustomerEmail,
		"total_amount":      order.TotalAmount,
		"currency":          order.Currency,
		"items_summary":     itemsSummary,
		"shipping_address":  shippingAddress,
		"payment_method_id": order.PaymentMethodID,
		"integration_id":    order.IntegrationID,
		"platform":          order.Platform,
		"timestamp":         time.Now().Unix(),
	}

	// Serializar a JSON
	payload, err := json.Marshal(event)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Msg("Error marshaling confirmation event")
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Publicar a RabbitMQ
	if err := p.rabbit.Publish(ctx, "orders.confirmation.requested", payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("queue", "orders.confirmation.requested").
			Msg("Error publishing confirmation event to RabbitMQ")
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	p.log.Info().
		Str("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Str("queue", "orders.confirmation.requested").
		Msg("Confirmation event published successfully")

	return nil
}

// buildItemsSummary construye un resumen de los items de la orden
func buildItemsSummary(order *entities.ProbabilityOrder) string {
	if len(order.OrderItems) == 0 {
		return "Sin items"
	}

	summary := ""
	for i, item := range order.OrderItems {
		if i > 0 {
			summary += ", "
		}
		summary += fmt.Sprintf("%dx %s", item.Quantity, item.ProductName)
	}

	return summary
}

// buildShippingAddress construye un resumen de la dirección de envío
func buildShippingAddress(order *entities.ProbabilityOrder) string {
	if order.ShippingStreet == "" && order.ShippingCity == "" {
		return "Sin dirección"
	}

	address := order.ShippingStreet
	if order.ShippingCity != "" {
		if address != "" {
			address += ", "
		}
		address += order.ShippingCity
	}
	if order.ShippingState != "" {
		if address != "" {
			address += ", "
		}
		address += order.ShippingState
	}

	return address
}

// ───────────────────────────────────────────
//
//	FUNCIONES HELPER
//
// ───────────────────────────────────────────

// publishToQueue publica un mensaje a una queue específica de RabbitMQ
func (p *OrderRabbitPublisher) publishToQueue(ctx context.Context, queue string, message *response.OrderEventMessage) error {
	// Serializar a JSON
	payload, err := json.Marshal(message)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", message.OrderID).
			Str("event_type", message.EventType).
			Msg("Error marshaling order event")
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Publicar a RabbitMQ
	if err := p.rabbit.Publish(ctx, queue, payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", message.OrderID).
			Str("event_type", message.EventType).
			Str("queue", queue).
			Msg("Error publishing order event to RabbitMQ")
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	p.log.Info().
		Str("order_id", message.OrderID).
		Str("event_type", message.EventType).
		Str("queue", queue).
		Msg("✅ Order event published to RabbitMQ")

	return nil
}

// getQueueForEventType retorna la queue correspondiente al tipo de evento
func (p *OrderRabbitPublisher) getQueueForEventType(eventType entities.OrderEventType) string {
	switch eventType {
	case entities.OrderEventTypeCreated:
		return "orders.events.created"
	case entities.OrderEventTypeUpdated:
		return "orders.events.updated"
	case entities.OrderEventTypeCancelled:
		return "orders.events.cancelled"
	case entities.OrderEventTypeStatusChanged:
		return "orders.events.status_changed"
	case entities.OrderEventTypeConfirmationRequested:
		return "orders.confirmation.requested"
	default:
		return "orders.events.generic"
	}
}
