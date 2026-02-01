package redis

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue/response"
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
func NewOrderEventPublisher(redisClient redisclient.IRedis, logger log.ILogger, channel string) ports.IOrderEventPublisher {
	return &OrderEventPublisher{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
	}
}

// PublishOrderEvent publica un evento de orden a Redis con snapshot completo
// Incluye toda la información de la orden para que consumidores no necesiten consultar BD
func (p *OrderEventPublisher) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	// Construir OrderSnapshot completo
	orderSnapshot := mappers.OrderToSnapshot(order)

	// Construir mensaje completo (OrderEventMessage)
	message := &response.OrderEventMessage{
		EventID:       event.ID,
		EventType:     string(event.Type),
		OrderID:       event.OrderID,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		Timestamp:     event.Timestamp,
		Order:         orderSnapshot, // ✅ SIEMPRE incluir snapshot completo
		Changes: map[string]interface{}{
			"previous_status": event.Data.PreviousStatus,
			"current_status":  event.Data.CurrentStatus,
			"platform":        event.Data.Platform,
		},
		Metadata: event.Metadata,
	}

	// Serializar mensaje completo a JSON
	eventJSON, err := json.Marshal(message)
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

	p.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Str("order_id", event.OrderID).
		Interface("business_id", event.BusinessID).
		Interface("integration_id", event.IntegrationID).
		Str("channel", p.channel).
		Str("customer_phone", order.CustomerPhone).
		Str("items_summary", orderSnapshot.ItemsSummary).
		Msg("✅ Evento de orden ENRIQUECIDO publicado a Redis Pub/Sub")

	return nil
}
