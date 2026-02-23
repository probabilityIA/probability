package events

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/primary"
	"github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/primary/handlers"
	eventsevents "github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/secondary/events"
	redisinfra "github.com/secamc93/probability/back/central/services/integrations/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// ─────────────────────────────────────────────────────────────────
// Re-exports públicos — alias de tipos internos del domain
// Permiten usar events.SyncOrderCreatedEvent{} sin importar internal/
// ─────────────────────────────────────────────────────────────────

// Interfaces
type IIntegrationEventService = domain.IIntegrationEventService

// Structs de eventos
type SyncOrderCreatedEvent  = domain.SyncOrderCreatedEvent
type SyncOrderUpdatedEvent  = domain.SyncOrderUpdatedEvent
type SyncOrderRejectedEvent = domain.SyncOrderRejectedEvent
type SyncStartedEvent       = domain.SyncStartedEvent
type SyncCompletedEvent     = domain.SyncCompletedEvent
type SyncFailedEvent        = domain.SyncFailedEvent
type SyncParams             = domain.SyncParams
type IntegrationEvent       = domain.IntegrationEvent
type IntegrationEventType   = domain.IntegrationEventType

// Constantes de tipos de evento
const (
	EventTypeSyncOrderCreated  = domain.IntegrationEventTypeSyncOrderCreated
	EventTypeSyncOrderUpdated  = domain.IntegrationEventTypeSyncOrderUpdated
	EventTypeSyncOrderRejected = domain.IntegrationEventTypeSyncOrderRejected
	EventTypeSyncStarted       = domain.IntegrationEventTypeSyncStarted
	EventTypeSyncCompleted     = domain.IntegrationEventTypeSyncCompleted
	EventTypeSyncFailed        = domain.IntegrationEventTypeSyncFailed
)

// eventServiceInstance es la instancia global del servicio de eventos (SSE in-process).
var eventServiceInstance domain.IIntegrationEventService

// redisPublisher publica eventos al canal Redis para que modules/events los consuma.
var redisPublisher *redisinfra.IntegrationEventRedisPublisher

const defaultIntegrationEventsChannel = redisclient.ChannelIntegrationsSyncOrders

// New inicializa el módulo de eventos de integraciones, registra las rutas
// y establece la instancia global para las funciones Publish* de este package.
func New(router *gin.RouterGroup, logger log.ILogger, redisClient redisclient.IRedis) domain.IIntegrationEventService {
	// 1. Init Event Manager (SSE y eventos en tiempo real — backward compat)
	eventManager := eventsevents.New(logger)

	// 2. Init Event Service
	eventService := app.NewIntegrationEventService(eventManager)

	// 3. Registrar instancia global (usada como fallback si Redis no está disponible)
	eventServiceInstance = eventService

	// 4. Init Redis Publisher (si Redis está disponible)
	if redisClient != nil {
		channel := defaultIntegrationEventsChannel
		redisPublisher = redisinfra.New(redisClient, channel, logger)
		redisClient.RegisterChannel(channel)
	}

	// 5. Init SSE Handler
	sseHandler := handlers.New(eventManager, logger)

	// 6. Init y registrar rutas
	primary.New(sseHandler).RegisterRoutes(router)

	return eventService
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return uuid.New().String()
}

// ─────────────────────────────────────────────────────────────────
// Helpers de publicación — construyen IntegrationEvent y lo
// publican a Redis. Si Redis no está disponible, usan el
// servicio in-memory como fallback.
// ─────────────────────────────────────────────────────────────────

func PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderCreatedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncOrderCreated,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncOrderCreated(ctx, integrationID, businessID, data) //nolint:errcheck
}

func PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderUpdatedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncOrderUpdated,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncOrderUpdated(ctx, integrationID, businessID, data) //nolint:errcheck
}

func PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderRejectedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncOrderRejected,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncOrderRejected(ctx, integrationID, businessID, data) //nolint:errcheck
}

func PublishSyncStarted(ctx context.Context, integrationID uint, businessID *uint, data SyncStartedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncStarted,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncStarted(ctx, integrationID, businessID, data) //nolint:errcheck
}

func PublishSyncCompleted(ctx context.Context, integrationID uint, businessID *uint, data SyncCompletedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncCompleted,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncCompleted(ctx, integrationID, businessID, data) //nolint:errcheck
}

func PublishSyncFailed(ctx context.Context, integrationID uint, businessID *uint, data SyncFailedEvent) {
	if redisPublisher != nil {
		redisPublisher.Publish(ctx, domain.IntegrationEvent{ //nolint:errcheck
			ID:            generateEventID(),
			Type:          domain.IntegrationEventTypeSyncFailed,
			IntegrationID: integrationID,
			BusinessID:    businessID,
			Timestamp:     time.Now(),
			Data:          data,
		})
		return
	}
	if eventServiceInstance == nil {
		return
	}
	eventServiceInstance.PublishSyncFailed(ctx, integrationID, businessID, data) //nolint:errcheck
}
