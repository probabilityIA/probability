package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	redisinfra "github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// ─────────────────────────────────────────────────────────────────
// Re-exports públicos — alias de tipos internos del domain
// Permiten usar events.SyncOrderCreatedEvent{} sin importar internal/
// ─────────────────────────────────────────────────────────────────

// Structs de eventos
type SyncOrderCreatedEvent = domain.SyncOrderCreatedEvent
type SyncOrderUpdatedEvent = domain.SyncOrderUpdatedEvent
type SyncOrderRejectedEvent = domain.SyncOrderRejectedEvent
type SyncStartedEvent = domain.SyncStartedEvent
type SyncCompletedEvent = domain.SyncCompletedEvent
type SyncFailedEvent = domain.SyncFailedEvent
type SyncParams = domain.SyncParams
type IntegrationEvent = domain.IntegrationEvent
type IntegrationEventType = domain.IntegrationEventType

// Constantes de tipos de evento
const (
	EventTypeSyncOrderCreated  = domain.IntegrationEventTypeSyncOrderCreated
	EventTypeSyncOrderUpdated  = domain.IntegrationEventTypeSyncOrderUpdated
	EventTypeSyncOrderRejected = domain.IntegrationEventTypeSyncOrderRejected
	EventTypeSyncStarted       = domain.IntegrationEventTypeSyncStarted
	EventTypeSyncCompleted     = domain.IntegrationEventTypeSyncCompleted
	EventTypeSyncFailed        = domain.IntegrationEventTypeSyncFailed
)

// redisPublisher publica eventos al canal Redis para que modules/events los consuma.
var redisPublisher *redisinfra.IntegrationEventRedisPublisher

const defaultIntegrationEventsChannel = redisclient.ChannelIntegrationsSyncOrders

// Init inicializa el publisher de eventos de integraciones.
// Solo configura el Redis publisher — la entrega SSE al frontend
// la maneja modules/events (centralizada).
func Init(logger log.ILogger, redisClient redisclient.IRedis) {
	if redisClient == nil {
		logger.Warn(context.Background()).Msg("Redis no disponible, integration events publisher no se inicializará")
		return
	}
	channel := defaultIntegrationEventsChannel
	redisPublisher = redisinfra.New(redisClient, channel, logger)
	redisClient.RegisterChannel(channel)
	logger.Info(context.Background()).
		Str("channel", channel).
		Msg("📤 Integration events Redis publisher inicializado")
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return uuid.New().String()
}

// ─────────────────────────────────────────────────────────────────
// Helpers de publicación — construyen IntegrationEvent y lo
// publican a Redis.
// ─────────────────────────────────────────────────────────────────

func PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderCreatedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderCreated,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}

func PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderUpdatedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderUpdated,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}

func PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderRejectedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncOrderRejected,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}

func PublishSyncStarted(ctx context.Context, integrationID uint, businessID *uint, data SyncStartedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncStarted,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}

func PublishSyncCompleted(ctx context.Context, integrationID uint, businessID *uint, data SyncCompletedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncCompleted,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}

func PublishSyncFailed(ctx context.Context, integrationID uint, businessID *uint, data SyncFailedEvent) {
	if redisPublisher == nil {
		return
	}
	redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
		ID:            generateEventID(),
		Type:          domain.IntegrationEventTypeSyncFailed,
		IntegrationID: integrationID,
		BusinessID:    businessID,
		Timestamp:     time.Now(),
		Data:          data,
	})
}
