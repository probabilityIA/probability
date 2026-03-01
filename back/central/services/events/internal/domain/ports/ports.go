package ports

import (
	"context"
	"net/http"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// ISSEPublisher define el puerto para manejar eventos SSE en tiempo real
type ISSEPublisher interface {
	AddConnection(businessID uint, filter *entities.SSEConnectionFilter, conn http.ResponseWriter) string
	RemoveConnection(connectionID string)
	PublishEvent(event entities.Event)
	GetConnectionCount(businessID uint) int
	GetConnectionInfo(businessID uint) map[string]interface{}
	GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []entities.Event
	HasRecentEvents(businessID uint) bool
	Stop()
}

// INotificationConfigCache define el puerto para leer configuraciones de notificación desde Redis cache
type INotificationConfigCache interface {
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.CachedNotificationConfig, error)
}

// IChannelPublisher define el puerto para publicar eventos a canales específicos (WhatsApp, Email)
type IChannelPublisher interface {
	PublishToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error
}

// IEventDispatcher define el puerto para el dispatcher de eventos (capa de aplicación)
type IEventDispatcher interface {
	HandleEvent(ctx context.Context, event entities.Event) error
}
