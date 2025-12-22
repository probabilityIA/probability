package domain

import (
	"context"
	"net/http"
)

// ───────────────────────────────────────────
//
//	IIntegrationEventPublisher - Puerto para manejar eventos de integraciones
//
// ───────────────────────────────────────────

// IIntegrationEventPublisher define el puerto para manejar eventos de integraciones en tiempo real
type IIntegrationEventPublisher interface {
	// Gestión de conexiones por business_id
	AddConnection(businessID uint, filter *IntegrationSSEFilter, conn http.ResponseWriter) string
	RemoveConnection(connectionID string)

	// Publicación de eventos
	PublishEvent(event IntegrationEvent)

	// Información del sistema
	GetConnectionCount(businessID uint) int
	GetConnectionInfo(businessID uint) map[string]interface{}

	// Historial/caché de eventos recientes por business_id
	GetRecentEventsByBusiness(businessID uint, sinceSeq int64) []IntegrationEvent
	HasRecentEvents(businessID uint) bool

	// Control del sistema
	Stop()
}

// ───────────────────────────────────────────
//
//	IIntegrationEventService - Puerto para el servicio de eventos
//
// ───────────────────────────────────────────

// IIntegrationEventService define el servicio para publicar eventos de integraciones
type IIntegrationEventService interface {
	// Publicar eventos de sincronización
	PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderCreatedEvent) error
	PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderUpdatedEvent) error
	PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data SyncOrderRejectedEvent) error
	PublishSyncStarted(ctx context.Context, integrationID uint, businessID *uint, data SyncStartedEvent) error
	PublishSyncCompleted(ctx context.Context, integrationID uint, businessID *uint, data SyncCompletedEvent) error
	PublishSyncFailed(ctx context.Context, integrationID uint, businessID *uint, data SyncFailedEvent) error

	// Gestión de estado de sincronización
	SetSyncState(ctx context.Context, integrationID uint, state SyncState) error
	GetSyncState(ctx context.Context, integrationID uint) (*SyncState, error)
	DeleteSyncState(ctx context.Context, integrationID uint) error
	IncrementSyncCounter(ctx context.Context, integrationID uint, counterType string) error
}
