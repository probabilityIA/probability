package domain

import (
	"context"
	"net/http"
)

// ───────────────────────────────────────────
//
//	IEventPublisher - Puerto para manejar eventos en tiempo real
//
// ───────────────────────────────────────────

// IEventPublisher define el puerto para manejar eventos en tiempo real
type IEventPublisher interface {
	// Gestión de conexiones por business_id
	AddConnection(businessID uint, filter *SSEConnectionFilter, conn http.ResponseWriter) string
	RemoveConnection(connectionID string)

	// Publicación de eventos
	PublishEvent(event Event)

	// Información del sistema
	GetConnectionCount(businessID uint) int
	GetConnectionInfo(businessID uint) map[string]interface{}

	// Historial/caché de eventos recientes por business_id
	GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []Event
	HasRecentEvents(businessID uint) bool

	// Control del sistema
	Stop()
}

// ───────────────────────────────────────────
//
//	INotificationConfigRepository - Puerto para repositorio de configuraciones
//
// ───────────────────────────────────────────

// INotificationConfigRepository define el repositorio de configuraciones de notificaciones
type INotificationConfigRepository interface {
	// CRUD Operations
	Create(ctx context.Context, config *NotificationConfig) error
	GetByID(ctx context.Context, id uint) (*NotificationConfig, error)
	GetByBusinessAndEventType(ctx context.Context, businessID uint, eventType string) (*NotificationConfig, error)
	GetByBusinessID(ctx context.Context, businessID uint) ([]NotificationConfig, error)
	List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]NotificationConfig, int64, error)
	Update(ctx context.Context, config *NotificationConfig) error
	Delete(ctx context.Context, id uint) error

	// Helpers
	IsEventTypeEnabled(ctx context.Context, businessID uint, eventType string) (bool, error)
}
