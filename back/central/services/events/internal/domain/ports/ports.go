package ports

import (
	"context"
	"net/http"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

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

type INotificationConfigCache interface {
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.CachedNotificationConfig, error)
}

type IChannelPublisher interface {
	PublishToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error
	PublishToEmail(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error
	PublishToPush(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error
}

type IEventDispatcher interface {
	HandleEvent(ctx context.Context, event entities.Event) error
}

type IAssistantAlertPublisher interface {
	PublishToAssistant(ctx context.Context, event entities.Event) error
}
