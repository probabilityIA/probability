package queue

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type orderEventMessage struct {
	EventType string         `json:"event_type"`
	OrderID   string         `json:"order_id"`
	Order     *orderSnapshot `json:"order"`
	Changes   map[string]any `json:"changes,omitempty"`
}

type orderSnapshot struct {
	Platform       string `json:"platform"`
	Status         string `json:"status"`
	ExternalID     string `json:"external_id"`
	IntegrationID  uint   `json:"integration_id"`
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}

type OrderStatusConsumer struct {
	queue  rabbitmq.IQueue
	uc     usecases.IJumpsellerUseCase
	logger log.ILogger
}

func NewOrderStatusConsumer(queue rabbitmq.IQueue, uc usecases.IJumpsellerUseCase, logger log.ILogger) *OrderStatusConsumer {
	return &OrderStatusConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("jumpseller.status_consumer"),
	}
}

func (c *OrderStatusConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, Jumpseller status push-back consumer disabled")
		return
	}

	if err := c.queue.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare orders events exchange")
		return
	}
	if err := c.queue.DeclareQueue(rabbitmq.QueueOrdersToJumpseller, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare jumpseller orders queue")
		return
	}
	if err := c.queue.BindQueue(rabbitmq.QueueOrdersToJumpseller, rabbitmq.ExchangeOrderEvents, ""); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to bind jumpseller orders queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueOrdersToJumpseller).Msg("Starting Jumpseller status push-back consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueOrdersToJumpseller, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Jumpseller status consumer stopped with error")
		}
	}()
}

func (c *OrderStatusConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg orderEventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal order event message")
		return
	}

	if msg.EventType != "order.status_changed" {
		return
	}
	if msg.Order == nil || !isJumpsellerPlatform(msg.Order.Platform) {
		return
	}
	if msg.Order.ExternalID == "" || msg.Order.IntegrationID == 0 {
		return
	}

	currentStatus, _ := msg.Changes["current_status"].(string)
	if currentStatus == "" {
		currentStatus = msg.Order.Status
	}
	if currentStatus == "" {
		return
	}

	integrationID := strconv.FormatUint(uint64(msg.Order.IntegrationID), 10)
	tracking := domain.UpdateOrderFields{
		TrackingNumber:  msg.Order.TrackingNumber,
		TrackingCompany: msg.Order.Carrier,
	}

	if err := c.uc.UpdateOrderStatus(ctx, integrationID, msg.Order.ExternalID, currentStatus, tracking); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("order_id", msg.OrderID).
			Str("external_id", msg.Order.ExternalID).
			Str("status", currentStatus).
			Msg("Failed to push order status to Jumpseller")
		return
	}

	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Str("external_id", msg.Order.ExternalID).
		Str("status", currentStatus).
		Msg("Pushed order status to Jumpseller")
}

func isJumpsellerPlatform(platform string) bool {
	return strings.Contains(strings.ToLower(platform), "jumpseller")
}
