package queue

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
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
	Platform string `json:"platform"`
	Status   string `json:"status"`
}

type OrderStatusConsumer struct {
	queue  rabbitmq.IQueue
	uc     usecases.IMeliUseCase
	repo   domain.IOrderLookupRepository
	logger log.ILogger
}

func NewOrderStatusConsumer(queue rabbitmq.IQueue, uc usecases.IMeliUseCase, repo domain.IOrderLookupRepository, logger log.ILogger) *OrderStatusConsumer {
	return &OrderStatusConsumer{
		queue:  queue,
		uc:     uc,
		repo:   repo,
		logger: logger.WithModule("meli.status_consumer"),
	}
}

func (c *OrderStatusConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, MercadoLibre status push-back consumer disabled")
		return
	}

	if err := c.queue.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare orders events exchange")
		return
	}
	if err := c.queue.DeclareQueue(rabbitmq.QueueOrdersToMeli, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare meli orders queue")
		return
	}
	if err := c.queue.BindQueue(rabbitmq.QueueOrdersToMeli, rabbitmq.ExchangeOrderEvents, ""); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to bind meli orders queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueOrdersToMeli).Msg("Starting MercadoLibre status push-back consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueOrdersToMeli, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("MercadoLibre status consumer stopped with error")
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
	if msg.Order != nil && msg.Order.Platform != "" && !isMercadoLibrePlatform(msg.Order.Platform) {
		return
	}

	currentStatus, _ := msg.Changes["current_status"].(string)
	if currentStatus == "" && msg.Order != nil {
		currentStatus = msg.Order.Status
	}
	if currentStatus == "" || msg.OrderID == "" {
		return
	}

	ref, err := c.repo.GetMeliShipmentByOrderID(ctx, msg.OrderID)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to lookup MercadoLibre shipment for order")
		return
	}
	if ref == nil {
		return
	}

	integrationID := strconv.FormatUint(uint64(ref.IntegrationID), 10)
	if err := c.uc.PushOrderStatus(ctx, integrationID, ref.ShipmentID, currentStatus); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("order_id", msg.OrderID).
			Int64("shipment_id", ref.ShipmentID).
			Str("status", currentStatus).
			Msg("Failed to push order status to MercadoLibre")
		return
	}

	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Int64("shipment_id", ref.ShipmentID).
		Str("status", currentStatus).
		Msg("Pushed order status to MercadoLibre")
}

func isMercadoLibrePlatform(platform string) bool {
	p := strings.ToLower(platform)
	return strings.Contains(p, "mercado") || p == "meli"
}
