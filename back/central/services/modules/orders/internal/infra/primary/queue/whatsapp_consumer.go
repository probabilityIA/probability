package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type WhatsAppConfirmedEvent struct {
	EventType   string `json:"event_type"`
	OrderNumber string `json:"order_number"`
	PhoneNumber string `json:"phone_number"`
	BusinessID  uint   `json:"business_id"`
	Source      string `json:"source"`
	Timestamp   int64  `json:"timestamp"`
}

type WhatsAppCancelledEvent struct {
	EventType          string `json:"event_type"`
	OrderNumber        string `json:"order_number"`
	CancellationReason string `json:"cancellation_reason"`
	PhoneNumber        string `json:"phone_number"`
	BusinessID         uint   `json:"business_id"`
	Source             string `json:"source"`
	Timestamp          int64  `json:"timestamp"`
}

type WhatsAppNoveltyEvent struct {
	EventType   string `json:"event_type"`
	OrderNumber string `json:"order_number"`
	NoveltyType string `json:"novelty_type"`
	PhoneNumber string `json:"phone_number"`
	BusinessID  uint   `json:"business_id"`
	Source      string `json:"source"`
	Timestamp   int64  `json:"timestamp"`
}

type WhatsAppConsumer struct {
	queue           rabbitmq.IQueue
	repository      ports.IRepository
	rabbitPublisher ports.IOrderRabbitPublisher
	log             log.ILogger
}

func NewWhatsAppConsumer(
	queue rabbitmq.IQueue,
	repository ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) *WhatsAppConsumer {
	return &WhatsAppConsumer{
		queue:           queue,
		repository:      repository,
		rabbitPublisher: rabbitPublisher,
		log:             logger,
	}
}

func (c *WhatsAppConsumer) Start(ctx context.Context) error {
	queues := []string{
		rabbitmq.QueueWhatsAppOrderConfirmed,
		rabbitmq.QueueWhatsAppOrderCancelled,
		rabbitmq.QueueWhatsAppOrderNovelty,
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

	go func() {
		if err := c.queue.Consume(ctx, rabbitmq.QueueWhatsAppOrderConfirmed, c.handleConfirmed); err != nil {
			c.log.Error().Err(err).Msg("Error consuming confirmed queue")
		}
	}()

	go func() {
		if err := c.queue.Consume(ctx, rabbitmq.QueueWhatsAppOrderCancelled, c.handleCancelled); err != nil {
			c.log.Error().Err(err).Msg("Error consuming cancelled queue")
		}
	}()

	go func() {
		if err := c.queue.Consume(ctx, rabbitmq.QueueWhatsAppOrderNovelty, c.handleNovelty); err != nil {
			c.log.Error().Err(err).Msg("Error consuming novelty queue")
		}
	}()

	return nil
}

func (c *WhatsAppConsumer) handleConfirmed(msg []byte) error {
	var event WhatsAppConfirmedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Warn().Err(err).Msg("Discarding confirmed event: payload malformado (ACK)")
		return nil
	}

	ctx := context.Background()

	c.log.Info().
		Str("order_number", event.OrderNumber).
		Str("phone_number", event.PhoneNumber).
		Msg("Processing order confirmation from WhatsApp")

	order, err := c.repository.GetOrderByOrderNumberAndBusiness(ctx, event.OrderNumber, event.BusinessID)
	if err != nil {
		if esPermanente(err) {
			c.log.Warn().
				Err(err).
				Str("order_number", event.OrderNumber).
				Uint("business_id", event.BusinessID).
				Msg("Discarding confirmed event: la orden no existe para ese negocio (ACK)")
			return nil
		}
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Uint("business_id", event.BusinessID).
			Msg("Error getting order for confirmation - will be retried")
		return err
	}

	previousStatus := "pending"
	if order.OrderStatus != nil {
		previousStatus = order.OrderStatus.Code
	}

	confirmed := true
	order.IsConfirmed = &confirmed

	processingStatusID, err := c.repository.GetOrderStatusIDByCode(ctx, "processing")
	if err != nil {
		c.log.Warn().Err(err).Msg("Error getting processing status ID, skipping status change")
	} else if processingStatusID != nil {
		order.StatusID = processingStatusID
	}

	if err := c.repository.UpdateOrder(ctx, order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order confirmation status")
		return err
	}

	if c.rabbitPublisher != nil {
		go func() {
			bgCtx := context.Background()
			if err := c.rabbitPublisher.PublishOrderStatusChanged(bgCtx, order, previousStatus, "processing"); err != nil {
				c.log.Error().Err(err).Msg("Error publishing confirmation event to fanout")
			}
		}()
	}

	c.log.Info().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Str("previous_status", previousStatus).
		Str("new_status", "processing").
		Msg("Order confirmed successfully via WhatsApp")

	return nil
}

func (c *WhatsAppConsumer) handleCancelled(msg []byte) error {
	var event WhatsAppCancelledEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Warn().Err(err).Msg("Discarding cancelled event: payload malformado (ACK)")
		return nil
	}

	c.log.Warn().
		Str("order_number", event.OrderNumber).
		Str("reason", event.CancellationReason).
		Msg("Processing order cancellation from WhatsApp")

	order, err := c.repository.GetOrderByOrderNumberAndBusiness(context.Background(), event.OrderNumber, event.BusinessID)
	if err != nil {
		if esPermanente(err) {
			c.log.Warn().
				Err(err).
				Str("order_number", event.OrderNumber).
				Uint("business_id", event.BusinessID).
				Msg("Discarding cancelled event: la orden no existe para ese negocio (ACK)")
			return nil
		}
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Msg("Error getting order for cancellation - will be retried")
		return err
	}

	confirmed := false
	noveltyText := fmt.Sprintf("Cancelación solicitada vía WhatsApp: %s (Teléfono: %s)", event.CancellationReason, event.PhoneNumber)

	if order.Novelty != nil && *order.Novelty != "" {
		noveltyText = *order.Novelty + " | " + noveltyText
	}

	order.IsConfirmed = &confirmed
	order.Novelty = &noveltyText

	if err := c.repository.UpdateOrder(context.Background(), order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order cancellation status")
		return err
	}

	if c.rabbitPublisher != nil {
		go func() {
			if err := c.rabbitPublisher.PublishOrderCancelled(context.Background(), order); err != nil {
				c.log.Error().Err(err).Msg("Error publishing cancellation event to fanout")
			}
		}()
	}

	c.log.Warn().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Msg("Order cancellation recorded - requires manual review")

	return nil
}

func (c *WhatsAppConsumer) handleNovelty(msg []byte) error {
	var event WhatsAppNoveltyEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Warn().Err(err).Msg("Discarding novelty event: payload malformado (ACK)")
		return nil
	}

	c.log.Info().
		Str("order_number", event.OrderNumber).
		Str("novelty_type", event.NoveltyType).
		Msg("Processing order novelty from WhatsApp")

	order, err := c.repository.GetOrderByOrderNumberAndBusiness(context.Background(), event.OrderNumber, event.BusinessID)
	if err != nil {
		if esPermanente(err) {
			c.log.Warn().
				Err(err).
				Str("order_number", event.OrderNumber).
				Uint("business_id", event.BusinessID).
				Msg("Discarding novelty event: la orden no existe para ese negocio (ACK)")
			return nil
		}
		c.log.Error().
			Err(err).
			Str("order_number", event.OrderNumber).
			Msg("Error getting order for novelty - will be retried")
		return err
	}

	var noveltyText string
	switch event.NoveltyType {
	case "change_address":
		noveltyText = fmt.Sprintf("El cliente solicita cambio de direccion via WhatsApp (Telefono: %s). La direccion NO se cambia sola: verificar con el cliente antes de generar la guia", event.PhoneNumber)
	case "change_products":
		noveltyText = fmt.Sprintf("El cliente solicita cambio de productos via WhatsApp (Telefono: %s)", event.PhoneNumber)
	case "change_payment":
		noveltyText = fmt.Sprintf("El cliente solicita cambio de medio de pago via WhatsApp (Telefono: %s)", event.PhoneNumber)
	default:
		noveltyText = fmt.Sprintf("Novedad via WhatsApp: %s (Telefono: %s)", event.NoveltyType, event.PhoneNumber)
	}

	if order.Novelty != nil && *order.Novelty != "" {
		noveltyText = *order.Novelty + " | " + noveltyText
	}

	order.Novelty = &noveltyText

	noConfirmada := false
	order.IsConfirmed = &noConfirmada

	if err := c.repository.UpdateOrder(context.Background(), order); err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", order.ID).
			Str("order_number", event.OrderNumber).
			Msg("Error updating order novelty")
		return err
	}

	if c.rabbitPublisher != nil {
		go func() {
			if err := c.rabbitPublisher.PublishOrderUpdated(context.Background(), order); err != nil {
				c.log.Error().Err(err).Msg("Error publishing novelty event to fanout")
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

func esPermanente(err error) bool {
	return errors.Is(err, domainerrors.ErrOrderNotFound) ||
		errors.Is(err, domainerrors.ErrOrderBusinessDeleted)
}
