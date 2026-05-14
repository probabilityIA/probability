package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type orderEventMessage struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	OrderID    string `json:"order_id"`
	BusinessID uint   `json:"business_id"`
}

type ProbabilityConsumer struct {
	queue rabbitmq.IQueue
	cache ports.IProbabilityCache
	log   log.ILogger
}

func NewProbabilityConsumer(queue rabbitmq.IQueue, cache ports.IProbabilityCache, logger log.ILogger) *ProbabilityConsumer {
	return &ProbabilityConsumer{
		queue: queue,
		cache: cache,
		log:   logger.WithModule("geozones.probability.consumer"),
	}
}

func (c *ProbabilityConsumer) Start(ctx context.Context) {
	if c.queue == nil || c.cache == nil {
		c.log.Warn(ctx).Msg("RabbitMQ or probability cache not available, consumer disabled")
		return
	}

	queueName := rabbitmq.QueueOrdersToGeozonesProbability

	if err := c.queue.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare order events exchange")
		return
	}
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error(ctx).Err(err).Str("queue", queueName).Msg("Failed to declare queue")
		return
	}
	if err := c.queue.BindQueue(queueName, rabbitmq.ExchangeOrderEvents, ""); err != nil {
		c.log.Error(ctx).Err(err).Str("queue", queueName).Msg("Failed to bind queue to fanout")
		return
	}

	c.log.Info(ctx).Str("queue", queueName).Msg("Starting geozones probability invalidation consumer")

	go func() {
		err := c.queue.Consume(ctx, queueName, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.log.Error(ctx).Err(err).Msg("Probability invalidation consumer stopped with error")
		}
	}()
}

func (c *ProbabilityConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg orderEventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.log.Warn(ctx).Err(err).Msg("Failed to unmarshal order event")
		return
	}
	if msg.OrderID == "" {
		return
	}
	if err := c.cache.InvalidateOrder(ctx, msg.BusinessID, msg.OrderID); err != nil {
		c.log.Warn(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to invalidate probability cache")
		return
	}
	c.log.Debug(ctx).
		Str("order_id", msg.OrderID).
		Str("event_type", msg.EventType).
		Uint("business_id", msg.BusinessID).
		Msg("probability cache invalidated by order event")
}
