package redis

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// OrderEventPublisher publica eventos de órdenes a Redis Pub/Sub
type OrderEventPublisher struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
}

// NewOrderEventPublisher crea un nuevo publicador de eventos de órdenes
func NewOrderEventPublisher(redisClient redisclient.IRedis, logger log.ILogger, channel string) domain.IOrderEventPublisher {
	return &OrderEventPublisher{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
	}
}

// PublishOrderEvent publica un evento de orden a Redis
func (p *OrderEventPublisher) PublishOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	// Serializar evento a JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Msg("Error al serializar evento de orden")
		return err
	}

	// Publicar a Redis
	if err := p.redisClient.Client(ctx).Publish(ctx, p.channel, eventJSON).Err(); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Str("channel", p.channel).
			Msg("Error al publicar evento de orden a Redis")
		return err
	}

	p.logger.Debug(ctx).
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Str("order_id", event.OrderID).
		Str("channel", p.channel).
		Msg("Evento de orden publicado a Redis")

	return nil
}

