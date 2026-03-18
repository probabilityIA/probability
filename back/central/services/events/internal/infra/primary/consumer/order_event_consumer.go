package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/events/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// orderEventMessage es una struct local para deserializar mensajes del fanout de órdenes.
// Replica la estructura de OrderEventMessage del módulo orders sin importarlo.
type orderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *orderSnapshot         `json:"order"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// orderSnapshot replica los campos necesarios del OrderSnapshot del módulo orders.
type orderSnapshot struct {
	ID              string  `json:"id"`
	OrderNumber     string  `json:"order_number"`
	InternalNumber  string  `json:"internal_number"`
	ExternalID      string  `json:"external_id"`
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email,omitempty"`
	CustomerPhone   string  `json:"customer_phone,omitempty"`
	Platform        string  `json:"platform"`
	IntegrationID   uint    `json:"integration_id"`
	BusinessName    string  `json:"business_name,omitempty"`
	ItemsSummary    string  `json:"items_summary,omitempty"`
	ShippingAddress string  `json:"shipping_address,omitempty"`
	OrderStatusID   *uint   `json:"order_status_id,omitempty"`
}

// OrderEventConsumer consume eventos de órdenes desde el fanout y los despacha al EventDispatcher
type OrderEventConsumer struct {
	rabbitMQ   rabbitmq.IQueue
	dispatcher ports.IEventDispatcher
	logger     log.ILogger
}

// NewOrderEventConsumer crea un nuevo consumer de eventos de órdenes
func NewOrderEventConsumer(
	rabbitMQ rabbitmq.IQueue,
	dispatcher ports.IEventDispatcher,
	logger log.ILogger,
) *OrderEventConsumer {
	return &OrderEventConsumer{
		rabbitMQ:   rabbitMQ,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

// Start inicia el consumer en background
func (c *OrderEventConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueOrdersToEvents).
		Msg("Iniciando consumer de eventos de órdenes (fanout → events dispatcher)")

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueOrdersToEvents, func(body []byte) error {
		return c.handleMessage(ctx, body)
	})
}

// handleMessage deserializa un OrderEventMessage y lo transforma a entities.Event
func (c *OrderEventConsumer) handleMessage(ctx context.Context, body []byte) error {
	var msg orderEventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		wrappedErr := fmt.Errorf("%w: %v", domainerrors.ErrDeserializeFailed, err)
		c.logger.Error(ctx).
			Err(wrappedErr).
			Str("body", string(body)).
			Msg("Error deserializando evento de orden desde fanout")
		return nil // No requeue mensajes malformados
	}

	// Extraer business_id e integration_id
	var businessID uint
	if msg.BusinessID != nil {
		businessID = *msg.BusinessID
	}
	var integrationID uint
	if msg.IntegrationID != nil {
		integrationID = *msg.IntegrationID
	}

	// Construir Data map con campos del snapshot y changes
	data := make(map[string]interface{})

	if msg.Order != nil {
		data["order_id"] = msg.Order.ID
		data["order_number"] = msg.Order.OrderNumber
		data["internal_number"] = msg.Order.InternalNumber
		data["external_id"] = msg.Order.ExternalID
		data["total_amount"] = msg.Order.TotalAmount
		data["currency"] = msg.Order.Currency
		data["customer_name"] = msg.Order.CustomerName
		data["customer_email"] = msg.Order.CustomerEmail
		data["customer_phone"] = msg.Order.CustomerPhone
		data["platform"] = msg.Order.Platform
		data["business_name"] = msg.Order.BusinessName
		data["items_summary"] = msg.Order.ItemsSummary
		data["shipping_address"] = msg.Order.ShippingAddress

		// Extraer order_status_id del snapshot (disponible en todos los eventos)
		if msg.Order.OrderStatusID != nil {
			data["order_status_id"] = *msg.Order.OrderStatusID
		}
	}

	// Extraer current_status de Changes (disponible en status_changed/updated)
	if msg.Changes != nil {
		if currentStatus, ok := msg.Changes["current_status"]; ok {
			data["current_status"] = currentStatus
		}
		if previousStatus, ok := msg.Changes["previous_status"]; ok {
			data["previous_status"] = previousStatus
		}
	}

	event := entities.Event{
		ID:            msg.EventID,
		Type:          msg.EventType,
		Category:      "order",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Timestamp:     msg.Timestamp,
		Data:          data,
		Metadata:      msg.Metadata,
	}

	c.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("business_id", event.BusinessID).
		Str("order_id", msg.OrderID).
		Msg("Evento de orden recibido desde fanout, despachando a EventDispatcher")

	return c.dispatcher.HandleEvent(ctx, event)
}
