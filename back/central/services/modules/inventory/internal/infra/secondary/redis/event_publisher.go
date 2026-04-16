package redis

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	sharedredis "github.com/secamc93/probability/back/central/shared/redis"
)

// EventPublisher publica eventos de inventario a Redis Pub/Sub para SSE
// y al dispatcher central via RabbitMQ (ExchangeEvents)
type EventPublisher struct {
	redis   sharedredis.IRedis
	logger  log.ILogger
	channel string
	queue   rabbitmq.IQueue
}

// NewEventPublisher crea un nuevo publisher de eventos de inventario
func NewEventPublisher(redisClient sharedredis.IRedis, logger log.ILogger, queue rabbitmq.IQueue) ports.IInventoryEventPublisher {
	if redisClient == nil {
		return &EventPublisher{redis: nil, logger: logger, queue: queue}
	}
	return &EventPublisher{
		redis:   redisClient,
		logger:  logger.WithModule("inventory.events"),
		channel: sharedredis.ChannelInventoryEvents,
		queue:   queue,
	}
}

// PublishInventoryEvent publica un evento de inventario al canal Redis
func (p *EventPublisher) PublishInventoryEvent(ctx context.Context, event ports.InventoryEvent) error {
	if p.redis == nil {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Error().Err(err).Msg("Failed to marshal inventory event")
		return err
	}

	if err := p.redis.Client(ctx).Publish(ctx, p.channel, body).Err(); err != nil {
		p.logger.Error().
			Err(err).
			Str("event_type", event.EventType).
			Str("order_id", event.OrderID).
			Msg("Failed to publish inventory event")
		return err
	}

	// Notificar al dispatcher central (SSE/WhatsApp/Email)
	if p.queue != nil {
		go func() {
			_ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
				Type:       event.EventType,
				Category:   "inventory",
				BusinessID: event.BusinessID,
				Data:       event.Data,
				Metadata: map[string]interface{}{
					"order_id":     event.OrderID,
					"warehouse_id": event.WarehouseID,
				},
			})
		}()
	}

	return nil
}
