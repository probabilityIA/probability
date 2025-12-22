package domain

import (
	"fmt"
	"time"
)

// SyncState representa el estado de una sincronización activa
type SyncState struct {
	IntegrationID   uint
	BusinessID      *uint
	IntegrationType string
	Status          SyncStatus
	StartedAt       time.Time
	Params          SyncParams
	TotalOrders     int
	CreatedOrders   int
	UpdatedOrders   int
	RejectedOrders  int
}

// SyncStatus representa el estado de una sincronización
type SyncStatus string

const (
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusCompleted  SyncStatus = "completed"
	SyncStatusFailed     SyncStatus = "failed"
)

// GetSyncStateKey genera la clave Redis para el estado de sincronización
func GetSyncStateKey(integrationID uint) string {
	return fmt.Sprintf("sync:orders:%d", integrationID)
}
