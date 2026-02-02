package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// OrderEventConsumer consume eventos de órdenes desde Redis Pub/Sub
type OrderEventConsumer struct {
	redisClient  redisclient.IRedis
	logger       log.ILogger
	channel      string
	scoreUseCase ports.IOrderScoreUseCase
}

// IOrderEventConsumer define la interfaz para consumir eventos de órdenes
type IOrderEventConsumer interface {
	Start(ctx context.Context) error
}

// NewOrderEventConsumer crea un nuevo consumidor de eventos de órdenes
func NewOrderEventConsumer(
	redisClient redisclient.IRedis,
	logger log.ILogger,
	channel string,
	scoreUseCase ports.IOrderScoreUseCase,
) IOrderEventConsumer {
	return &OrderEventConsumer{
		redisClient:  redisClient,
		logger:       logger,
		channel:      channel,
		scoreUseCase: scoreUseCase,
	}
}

// Start inicia el consumidor de eventos
func (c *OrderEventConsumer) Start(ctx context.Context) error {
	pubsub := c.redisClient.Client(ctx).Subscribe(ctx, c.channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info(ctx).Msg("Order event consumer stopped")
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			c.logger.Debug(ctx).
				Str("channel", c.channel).
				Str("payload", msg.Payload).
				Msg("Event received from Redis")

			// Deserializar evento
			var event entities.OrderEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				c.logger.Error(ctx).
					Err(err).
					Str("payload", msg.Payload).
					Msg("Error al deserializar evento de orden")
				continue
			}

			// Procesar evento según su tipo
			if err := c.handleEvent(ctx, &event); err != nil {
				c.logger.Error(ctx).
					Err(err).
					Str("event_id", event.ID).
					Str("event_type", string(event.Type)).
					Str("order_id", event.OrderID).
					Msg("Error al procesar evento de orden")
			}
		}
	}
}

// handleEvent procesa un evento según su tipo
func (c *OrderEventConsumer) handleEvent(ctx context.Context, event *entities.OrderEvent) error {
	switch event.Type {
	case entities.OrderEventTypeScoreCalculationRequested:
		fmt.Printf("[OrderEventConsumer] EVENTO RECIBIDO: order.score_calculation_requested para orden %s\n", event.OrderID)
		return c.scoreUseCase.CalculateAndUpdateOrderScore(ctx, event.OrderID)
	default:
		// Ignorar otros tipos de eventos
		return nil
	}
}
