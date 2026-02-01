package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// WhatsAppConfirmedEvent representa el evento de confirmación desde WhatsApp
type WhatsAppConfirmedEvent struct {
	EventType   string `json:"event_type"`
	OrderNumber string `json:"order_number"`
	PhoneNumber string `json:"phone_number"`
	BusinessID  uint   `json:"business_id"`
	Source      string `json:"source"`
	Timestamp   int64  `json:"timestamp"`
}

// WhatsAppCancelledEvent representa el evento de cancelación desde WhatsApp
type WhatsAppCancelledEvent struct {
	EventType          string `json:"event_type"`
	OrderNumber        string `json:"order_number"`
	CancellationReason string `json:"cancellation_reason"`
	PhoneNumber        string `json:"phone_number"`
	BusinessID         uint   `json:"business_id"`
	Source             string `json:"source"`
	Timestamp          int64  `json:"timestamp"`
}

// WhatsAppNoveltyEvent representa el evento de novedad desde WhatsApp
type WhatsAppNoveltyEvent struct {
	EventType   string `json:"event_type"`
	OrderNumber string `json:"order_number"`
	NoveltyType string `json:"novelty_type"`
	PhoneNumber string `json:"phone_number"`
	BusinessID  uint   `json:"business_id"`
	Source      string `json:"source"`
	Timestamp   int64  `json:"timestamp"`
}

// WhatsAppConsumer consume eventos de WhatsApp y actualiza órdenes
type WhatsAppConsumer struct {
	queue          rabbitmq.IQueue
	repository     ports.IRepository
	updateOrderUC  ports.IOrderUseCase // Interfaz en lugar de tipo concreto
	eventPublisher ports.IOrderEventPublisher
	log            log.ILogger
}

// NewWhatsAppConsumer crea un nuevo consumidor de eventos de WhatsApp
func NewWhatsAppConsumer(
	queue rabbitmq.IQueue,
	updateOrderUC ports.IOrderUseCase, // Interfaz en lugar de tipo concreto
	repository ports.IRepository,
	eventPublisher ports.IOrderEventPublisher,
	logger log.ILogger,
) *WhatsAppConsumer {
	return &WhatsAppConsumer{
		queue:          queue,
		updateOrderUC:  updateOrderUC,
		repository:     repository,
		eventPublisher: eventPublisher,
		log:            logger,
	}
}

// Start inicia el consumidor de eventos de WhatsApp
func (c *WhatsAppConsumer) Start(ctx context.Context) error {
	c.log.Info().Msg("Starting WhatsApp consumer for order confirmations")

	// Declarar colas
	queues := []string{
		"orders.whatsapp.confirmed",
		"orders.whatsapp.cancelled",
		"orders.whatsapp.novelty",
	}

	for _, queueName := range queues {
		if err := c.queue.DeclareQueue(queueName, true); err != nil {
			c.log.Error().
				Err(err).
				Str("queue", queueName).
				Msg("Error declaring queue")
			return err
		}
	}

	// Consumir de múltiples colas
	go func() {
		if err := c.queue.Consume(ctx, "orders.whatsapp.confirmed", c.handleConfirmed); err != nil {
			c.log.Error().Err(err).Msg("Error consuming confirmed queue")
		}
	}()

	go func() {
		if err := c.queue.Consume(ctx, "orders.whatsapp.cancelled", c.handleCancelled); err != nil {
			c.log.Error().Err(err).Msg("Error consuming cancelled queue")
		}
	}()

	go func() {
		if err := c.queue.Consume(ctx, "orders.whatsapp.novelty", c.handleNovelty); err != nil {
			c.log.Error().Err(err).Msg("Error consuming novelty queue")
		}
	}()

	c.log.Info().Msg("WhatsApp consumer started successfully")
	return nil
}

// handleConfirmed procesa eventos de confirmación
func (c *WhatsAppConsumer) handleConfirmed(msg []byte) error {
	var event WhatsAppConfirmedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Error().Err(err).Msg("Error unmarshaling confirmed event")
		return err
	}

	c.log.Info().
		Str("order_number", event.OrderNumber).
		Str("phone_number", event.PhoneNumber).
		Msg("Processing order confirmation from WhatsApp")

	// Buscar orden por order_number
	order, err := c.repository.GetOrderByOrderNumber(context.Background(), event.OrderNumber)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Msg("Error getting order for confirmation")
		return err
	}

	// Actualizar IsConfirmed = true
	confirmed := true
	order.IsConfirmed = &confirmed

	// Guardar cambios
	if err := c.repository.UpdateOrder(context.Background(), order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order confirmation status")
		return err
	}

	c.log.Info().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Msg("Order confirmed successfully")

	return nil
}

// handleCancelled procesa eventos de cancelación
func (c *WhatsAppConsumer) handleCancelled(msg []byte) error {
	var event WhatsAppCancelledEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Error().Err(err).Msg("Error unmarshaling cancelled event")
		return err
	}

	c.log.Warn().
		Str("order_number", event.OrderNumber).
		Str("reason", event.CancellationReason).
		Msg("Processing order cancellation from WhatsApp")

	// Buscar orden por order_number
	order, err := c.repository.GetOrderByOrderNumber(context.Background(), event.OrderNumber)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Msg("Error getting order for cancellation")
		return err
	}

	// Marcar IsConfirmed = false y guardar motivo en Novelty
	confirmed := false
	noveltyText := fmt.Sprintf("Cancelación solicitada vía WhatsApp: %s (Teléfono: %s)", event.CancellationReason, event.PhoneNumber)

	// Si ya existe novedad previa, concatenar
	if order.Novelty != nil && *order.Novelty != "" {
		noveltyText = *order.Novelty + " | " + noveltyText
	}

	order.IsConfirmed = &confirmed
	order.Novelty = &noveltyText

	// Guardar cambios
	if err := c.repository.UpdateOrder(context.Background(), order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order cancellation status")
		return err
	}

	// Publicar evento para notificar al equipo
	if c.eventPublisher != nil {
		go func() {
			orderEvent := entities.NewOrderEvent(
				entities.OrderEventTypeCancelled,
				order.ID,
				entities.OrderEventData{
					OrderNumber: event.OrderNumber,
					Extra: map[string]interface{}{
						"cancellation_source":    "whatsapp",
						"cancellation_reason":    event.CancellationReason,
						"requires_manual_review": true,
					},
				},
			)
			if err := c.eventPublisher.PublishOrderEvent(context.Background(), orderEvent, order); err != nil {
				c.log.Error().Err(err).Msg("Error publishing cancellation event")
			}
		}()
	}

	c.log.Warn().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Msg("Order cancellation recorded - requires manual review")

	return nil
}

// handleNovelty procesa eventos de novedades
func (c *WhatsAppConsumer) handleNovelty(msg []byte) error {
	var event WhatsAppNoveltyEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Error().Err(err).Msg("Error unmarshaling novelty event")
		return err
	}

	c.log.Info().
		Str("order_number", event.OrderNumber).
		Str("novelty_type", event.NoveltyType).
		Msg("Processing order novelty from WhatsApp")

	// Buscar orden por order_number
	order, err := c.repository.GetOrderByOrderNumber(context.Background(), event.OrderNumber)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Msg("Error getting order for novelty")
		return err
	}

	// Construir texto de novedad según tipo
	var noveltyText string
	switch event.NoveltyType {
	case "change_address":
		noveltyText = fmt.Sprintf("Solicitud de cambio de dirección vía WhatsApp (Teléfono: %s)", event.PhoneNumber)
	case "change_products":
		noveltyText = fmt.Sprintf("Solicitud de cambio de productos vía WhatsApp (Teléfono: %s)", event.PhoneNumber)
	case "change_payment":
		noveltyText = fmt.Sprintf("Solicitud de cambio de método de pago vía WhatsApp (Teléfono: %s)", event.PhoneNumber)
	default:
		noveltyText = fmt.Sprintf("Novedad vía WhatsApp: %s (Teléfono: %s)", event.NoveltyType, event.PhoneNumber)
	}

	// Si ya existe novedad previa, concatenar
	if order.Novelty != nil && *order.Novelty != "" {
		noveltyText = *order.Novelty + " | " + noveltyText
	}

	order.Novelty = &noveltyText

	// Guardar cambios
	if err := c.repository.UpdateOrder(context.Background(), order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order novelty")
		return err
	}

	// Publicar evento para notificar al equipo
	if c.eventPublisher != nil {
		go func() {
			orderEvent := entities.NewOrderEvent(
				entities.OrderEventTypeUpdated,
				order.ID,
				entities.OrderEventData{
					OrderNumber: event.OrderNumber,
					Extra: map[string]interface{}{
						"novelty_source":         "whatsapp",
						"novelty_type":           event.NoveltyType,
						"requires_manual_action": true,
					},
				},
			)
			if err := c.eventPublisher.PublishOrderEvent(context.Background(), orderEvent, order); err != nil {
				c.log.Error().Err(err).Msg("Error publishing novelty event")
			}
		}()
	}

	c.log.Info().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Str("novelty_type", event.NoveltyType).
		Msg("Order novelty recorded successfully")

	return nil
}
