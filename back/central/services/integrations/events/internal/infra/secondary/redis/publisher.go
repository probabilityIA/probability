package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// integrationEventPayload es la estructura de serialización JSON.
// El domain struct no tiene tags, así que usamos este struct para serializar.
type integrationEventPayload struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"event_type"`
	IntegrationID uint                   `json:"integration_id"`
	BusinessID    *uint                  `json:"business_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          interface{}            `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// IntegrationEventRedisPublisher publica eventos de integración a Redis Pub/Sub
type IntegrationEventRedisPublisher struct {
	redisClient redisclient.IRedis
	channel     string
	logger      log.ILogger
}

// New crea un nuevo publisher de eventos de integración a Redis
func New(
	redisClient redisclient.IRedis,
	channel string,
	logger log.ILogger,
) *IntegrationEventRedisPublisher {
	return &IntegrationEventRedisPublisher{
		redisClient: redisClient,
		channel:     channel,
		logger:      logger,
	}
}

// Publish serializa y publica un IntegrationEvent al canal Redis
func (p *IntegrationEventRedisPublisher) Publish(ctx context.Context, event domain.IntegrationEvent) error {
	payload := integrationEventPayload{
		ID:            event.ID,
		Type:          string(event.Type),
		IntegrationID: event.IntegrationID,
		BusinessID:    event.BusinessID,
		Timestamp:     event.Timestamp,
		Data:          event.Data,
		Metadata:      event.Metadata,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Msg("Error serializando integration event para Redis")
		return fmt.Errorf("error serializando integration event: %w", err)
	}

	client := p.redisClient.Client(ctx)
	if client == nil {
		return fmt.Errorf("redis client no disponible")
	}

	if err := client.Publish(ctx, p.channel, jsonBytes).Err(); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Str("channel", p.channel).
			Msg("Error publicando integration event a Redis")
		return fmt.Errorf("error publicando integration event a Redis: %w", err)
	}

	p.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Uint("integration_id", event.IntegrationID).
		Str("channel", p.channel).
		Msg("📤 Integration event publicado a Redis")

	return nil
}
