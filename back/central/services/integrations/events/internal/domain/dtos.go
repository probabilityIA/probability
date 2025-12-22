package domain

import (
	"time"
)

// ───────────────────────────────────────────
//
//	EVENT TYPES
//
// ───────────────────────────────────────────

// IntegrationEventType define los tipos de eventos de integraciones
type IntegrationEventType string

const (
	// Eventos de sincronización de órdenes
	IntegrationEventTypeSyncOrderCreated  IntegrationEventType = "integration.sync.order.created"
	IntegrationEventTypeSyncOrderUpdated  IntegrationEventType = "integration.sync.order.updated"
	IntegrationEventTypeSyncOrderRejected IntegrationEventType = "integration.sync.order.rejected"
	IntegrationEventTypeSyncStarted       IntegrationEventType = "integration.sync.started"
	IntegrationEventTypeSyncCompleted     IntegrationEventType = "integration.sync.completed"
	IntegrationEventTypeSyncFailed        IntegrationEventType = "integration.sync.failed"
)

// ───────────────────────────────────────────
//
//	EVENT STRUCTURES
//
// ───────────────────────────────────────────

// IntegrationEvent es la estructura base para todos los eventos de integraciones
type IntegrationEvent struct {
	ID            string
	Type          IntegrationEventType
	IntegrationID uint
	BusinessID    *uint
	Timestamp     time.Time
	Data          interface{}
	Metadata      map[string]interface{}
}

// SyncOrderCreatedEvent representa un evento de orden creada exitosamente durante la sincronización
type SyncOrderCreatedEvent struct {
	OrderID       string
	OrderNumber   string
	ExternalID    string
	Platform      string
	CustomerEmail string
	TotalAmount   *float64
	Currency      string
	Status        string
	CreatedAt     time.Time
	SyncedAt      time.Time
}

// SyncOrderUpdatedEvent representa un evento de orden actualizada exitosamente durante la sincronización
type SyncOrderUpdatedEvent struct {
	OrderID       string
	OrderNumber   string
	ExternalID    string
	Platform      string
	CustomerEmail string
	TotalAmount   *float64
	Currency      string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// SyncOrderRejectedEvent representa un evento de orden rechazada durante la sincronización
type SyncOrderRejectedEvent struct {
	OrderID     string
	OrderNumber string
	ExternalID  string
	Platform    string
	Reason      string
	Error       string
	RejectedAt  time.Time
}

// SyncStartedEvent representa el inicio de una sincronización
type SyncStartedEvent struct {
	IntegrationID   uint
	IntegrationType string
	Params          SyncParams
	StartedAt       time.Time
}

// SyncCompletedEvent representa la finalización exitosa de una sincronización
type SyncCompletedEvent struct {
	IntegrationID   uint
	IntegrationType string
	TotalOrders     int
	CreatedOrders   int
	UpdatedOrders   int
	RejectedOrders  int
	Duration        time.Duration
	CompletedAt     time.Time
}

// SyncFailedEvent representa el fallo de una sincronización
type SyncFailedEvent struct {
	IntegrationID   uint
	IntegrationType string
	Error           string
	FailedAt        time.Time
}

// SyncParams contiene los parámetros de sincronización
type SyncParams struct {
	CreatedAtMin      *time.Time
	CreatedAtMax      *time.Time
	Status            string
	FinancialStatus   string
	FulfillmentStatus string
}
