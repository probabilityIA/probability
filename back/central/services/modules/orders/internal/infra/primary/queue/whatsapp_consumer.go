package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
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
	queue             rabbitmq.IQueue
	repository        ports.IRepository
	rabbitPublisher   ports.IOrderRabbitPublisher
	statusUseCase     ports.IOrderStatusUseCase
	templateRequester ports.IWhatsAppTemplateRequester
	log               log.ILogger
}

func NewWhatsAppConsumer(
	queue rabbitmq.IQueue,
	repository ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	statusUseCase ports.IOrderStatusUseCase,
	templateRequester ports.IWhatsAppTemplateRequester,
	logger log.ILogger,
) *WhatsAppConsumer {
	return &WhatsAppConsumer{
		queue:             queue,
		repository:        repository,
		rabbitPublisher:   rabbitPublisher,
		statusUseCase:     statusUseCase,
		templateRequester: templateRequester,
		log:               logger,
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
	ctx := context.Background()

	var event WhatsAppCancelledEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		c.log.Warn().Err(err).Msg("Discarding cancelled event: payload malformado (ACK)")
		return nil
	}

	order, err := c.repository.GetOrderByOrderNumberAndBusiness(ctx, event.OrderNumber, event.BusinessID)
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

	current := entities.OrderStatus(order.Status)
	if current == entities.OrderStatusCancelled || current == entities.OrderStatusCancelRequested || current.IsTerminal() {
		c.log.Warn().
			Str("order_number", event.OrderNumber).
			Str("status", order.Status).
			Msg("Discarding cancelled event: la orden ya esta cancelada o con solicitud de cancelacion (ACK)")
		return nil
	}

	withGuide := HasActiveGuide(order)
	target := entities.OrderStatusCancelled
	template := templateOrderCancelled
	reason := "El cliente cancel\u00f3 el pedido por WhatsApp"
	if withGuide {
		target = entities.OrderStatusCancelRequested
		template = templateCancelRequestReceived
		reason = "El cliente solicit\u00f3 cancelar por WhatsApp y la orden ya tiene gu\u00eda: revisar si la transportadora ya la recogi\u00f3"
	}
	if strings.TrimSpace(event.CancellationReason) != "" {
		reason += ". Motivo: " + strings.TrimSpace(event.CancellationReason)
	}

	confirmed := false
	noveltyText := fmt.Sprintf("%s (Tel\u00e9fono: %s)", reason, event.PhoneNumber)
	if order.Novelty != nil && *order.Novelty != "" {
		noveltyText = *order.Novelty + " | " + noveltyText
	}
	order.IsConfirmed = &confirmed
	order.Novelty = &noveltyText
	if err := c.repository.UpdateOrder(ctx, order); err != nil {
		c.log.Error().Err(err).Str("order_id", order.ID).Msg("Error saving cancellation novelty - will be retried")
		return err
	}

	if c.statusUseCase != nil {
		_, err := c.statusUseCase.ChangeStatus(ctx, order.ID, &dtos.ChangeStatusRequest{
			Status:   string(target),
			UserName: "WhatsApp",
			Metadata: map[string]interface{}{"reason": reason, "source": "whatsapp"},
		})
		if err != nil {
			if errors.Is(err, domainerrors.ErrInvalidStatusTransition) || errors.Is(err, domainerrors.ErrOrderInTerminalState) {
				c.log.Warn().Err(err).Str("order_number", event.OrderNumber).
					Msg("Discarding cancelled event: el estado actual no admite el cambio (ACK)")
				return nil
			}
			c.log.Error().Err(err).Str("order_id", order.ID).Msg("Error changing status for WhatsApp cancellation - will be retried")
			return err
		}
	}

	if target == entities.OrderStatusCancelled && c.rabbitPublisher != nil {
		if err := c.rabbitPublisher.PublishOrderCancelled(ctx, order); err != nil {
			c.log.Error().Err(err).Msg("Error publishing cancellation event to fanout")
		}
	}

	if c.templateRequester != nil && event.PhoneNumber != "" {
		if err := c.templateRequester.RequestTemplate(ctx, event.BusinessID, event.PhoneNumber, template, []string{order.OrderNumber}); err != nil {
			c.log.Error().Err(err).Str("order_number", event.OrderNumber).Str("template", template).
				Msg("Error requesting WhatsApp cancellation template")
		}
	}

	c.log.Info().
		Str("order_id", order.ID).
		Str("order_number", event.OrderNumber).
		Str("new_status", string(target)).
		Bool("with_guide", withGuide).
		Msg("Cancelacion por WhatsApp procesada")

	return nil
}

const (
	templateOrderCancelled        = "pedido_cancelado"
	templateCancelRequestReceived = "solicitud_cancelacion_recibida"
)

func HasActiveGuide(order *entities.ProbabilityOrder) bool {
	if order == nil {
		return false
	}
	if len(order.Shipments) > 0 {
		shipment := order.Shipments[0]
		if shipment.DeletedAt != nil || strings.EqualFold(shipment.Status, "cancelled") {
			return false
		}
		return nonEmpty(shipment.TrackingNumber) || nonEmpty(shipment.GuideID)
	}
	return nonEmpty(order.TrackingNumber) || nonEmpty(order.GuideID)
}

func nonEmpty(value *string) bool {
	return value != nil && strings.TrimSpace(*value) != ""
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
