package domain

import "time"

// ───────────────────────────────────────────
//
//	INTEGRATION EVENT TYPES
//
// ───────────────────────────────────────────

// IntegrationEventType define los tipos de eventos de integraciones
type IntegrationEventType string

const (
	IntegrationEventTypeSyncOrderCreated  IntegrationEventType = "integration.sync.order.created"
	IntegrationEventTypeSyncOrderUpdated  IntegrationEventType = "integration.sync.order.updated"
	IntegrationEventTypeSyncOrderRejected IntegrationEventType = "integration.sync.order.rejected"
	IntegrationEventTypeSyncStarted       IntegrationEventType = "integration.sync.started"
	IntegrationEventTypeSyncCompleted     IntegrationEventType = "integration.sync.completed"
	IntegrationEventTypeSyncFailed        IntegrationEventType = "integration.sync.failed"
)

// IsValid verifica si el tipo de evento de integración es válido
func (t IntegrationEventType) IsValid() bool {
	switch t {
	case IntegrationEventTypeSyncOrderCreated,
		IntegrationEventTypeSyncOrderUpdated,
		IntegrationEventTypeSyncOrderRejected,
		IntegrationEventTypeSyncStarted,
		IntegrationEventTypeSyncCompleted,
		IntegrationEventTypeSyncFailed:
		return true
	}
	return false
}

// ───────────────────────────────────────────
//
//	INTEGRATION EVENT STRUCTURES
//
// ───────────────────────────────────────────

// IntegrationEvent representa un evento de integración recibido desde Redis.
// Los tags JSON coinciden con el payload publicado por integrations/events.
type IntegrationEvent struct {
	ID            string                 `json:"id"`
	Type          IntegrationEventType   `json:"event_type"`
	IntegrationID uint                   `json:"integration_id"`
	BusinessID    *uint                  `json:"business_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
