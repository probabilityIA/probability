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

//
//	MÉTODOS DE PUBLICACIÓN DE EVENTOS
//

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
	return p.publishToQueue(ctx, rabbitmq.RoutingKeyOrderCreated, message)
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
	return p.publishToQueue(ctx, rabbitmq.RoutingKeyOrderUpdated, message)
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
	return p.publishToQueue(ctx, rabbitmq.RoutingKeyOrderCancelled, message)
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
		Changes: map[string]any{
			"previous_status": previousStatus,
			"current_status":  currentStatus,
		},
		Order: mappers.OrderToSnapshot(order),
	}
	return p.publishToQueue(ctx, rabbitmq.RoutingKeyOrderStatusChanged, message)
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
		Changes: map[string]any{
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
	event := map[string]any{
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
	if err := p.rabbit.Publish(ctx, rabbitmq.QueueOrdersConfirmationRequested, payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("queue", rabbitmq.QueueOrdersConfirmationRequested).
			Msg("Error publishing confirmation event to RabbitMQ")
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	p.log.Info().
		Str("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Str("queue", rabbitmq.QueueOrdersConfirmationRequested).
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

// PublishGuideNotificationRequested publica un evento cuando se solicita enviar la guia por WhatsApp
func (p *OrderRabbitPublisher) PublishGuideNotificationRequested(ctx context.Context, order *entities.ProbabilityOrder) error {
	trackingNumber := ""
	if order.TrackingNumber != nil {
		trackingNumber = *order.TrackingNumber
	}

	event := map[string]any{
		"event_type":      "order.guide_notification_requested",
		"order_id":        order.ID,
		"order_number":    order.OrderNumber,
		"business_id":     order.BusinessID,
		"customer_name":   order.CustomerName,
		"customer_phone":  order.CustomerPhone,
		"tracking_number": trackingNumber,
		"integration_id":  order.IntegrationID,
		"platform":        order.Platform,
		"timestamp":       time.Now().Unix(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Msg("Error marshaling guide notification event")
		return fmt.Errorf("error marshaling event: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueShipmentsWhatsAppGuideNotification, payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("queue", rabbitmq.QueueShipmentsWhatsAppGuideNotification).
			Msg("Error publishing guide notification event to RabbitMQ")
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	p.log.Info().
		Str("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Str("tracking_number", trackingNumber).
		Msg("Guide notification event published successfully")

	return nil
}

//
//	FUNCIONES HELPER
//

// publishToQueue publica un mensaje a una queue específica de RabbitMQ
func (p *OrderRabbitPublisher) publishToQueue(ctx context.Context, _ string, message *response.OrderEventMessage) error {
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

	// Publicar al exchange fanout que distribuye a las 4 colas bindeadas
	routingKey := "" // Vacío porque es fanout (envía a todas las colas bindeadas)

	p.log.Debug().
		Str("order_id", message.OrderID).
		Str("event_type", message.EventType).
		Str("exchange", rabbitmq.ExchangeOrderEvents).
		Int("payload_size", len(payload)).
		Msg("📤 Publishing order event to exchange (fanout distribution)")

	if err := p.rabbit.PublishToExchange(ctx, rabbitmq.ExchangeOrderEvents, routingKey, payload); err != nil {
		p.log.Error().
			Err(err).
			Str("order_id", message.OrderID).
			Str("event_type", message.EventType).
			Str("exchange", rabbitmq.ExchangeOrderEvents).
			Msg("Error publishing order event to exchange")
		return fmt.Errorf("error publishing to exchange: %w", err)
	}

	p.log.Info().
		Str("order_id", message.OrderID).
		Str("event_type", message.EventType).
		Str("exchange", rabbitmq.ExchangeOrderEvents).
		Msg("Order event published to fanout (invoicing, score, inventory, events, customers)")

	return nil
}

// getQueueForEventType retorna la queue correspondiente al tipo de evento
func (p *OrderRabbitPublisher) getQueueForEventType(eventType entities.OrderEventType) string {
	switch eventType {
	case entities.OrderEventTypeCreated:
		return rabbitmq.RoutingKeyOrderCreated
	case entities.OrderEventTypeUpdated:
		return rabbitmq.RoutingKeyOrderUpdated
	case entities.OrderEventTypeCancelled:
		return rabbitmq.RoutingKeyOrderCancelled
	case entities.OrderEventTypeStatusChanged:
		return rabbitmq.RoutingKeyOrderStatusChanged
	case entities.OrderEventTypeConfirmationRequested:
		return rabbitmq.QueueOrdersConfirmationRequested
	default:
		return rabbitmq.RoutingKeyOrderGeneric
	}
}
