package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	NOTIFICATION CONFIGURATION DOMAIN
//
// ───────────────────────────────────────────

// NotificationConfig representa la configuración de notificaciones en el dominio
type NotificationConfig struct {
	ID          uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID  uint       `json:"business_id"`
	EventType   string     `json:"event_type"`
	Enabled     bool       `json:"enabled"`
	Channels    datatypes.JSON `json:"channels"`
	Filters     datatypes.JSON `json:"filters,omitempty"`
	Description string     `json:"description,omitempty"`
}

// ───────────────────────────────────────────
//
//	NOTIFICATION CONFIG DTOs
//
// ───────────────────────────────────────────

// CreateNotificationConfigRequest representa la solicitud para crear una configuración
type CreateNotificationConfigRequest struct {
	BusinessID  uint           `json:"business_id" binding:"required"`
	EventType   string         `json:"event_type" binding:"required"`
	Enabled     bool           `json:"enabled"`
	Channels    datatypes.JSON `json:"channels"`
	Filters     datatypes.JSON `json:"filters,omitempty"`
	Description string         `json:"description,omitempty"`
}

// UpdateNotificationConfigRequest representa la solicitud para actualizar una configuración
type UpdateNotificationConfigRequest struct {
	Enabled     *bool          `json:"enabled"`
	Channels    datatypes.JSON `json:"channels,omitempty"`
	Filters     datatypes.JSON `json:"filters,omitempty"`
	Description *string        `json:"description,omitempty"`
}

// NotificationConfigResponse representa la respuesta de una configuración
type NotificationConfigResponse struct {
	ID          uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	BusinessID  uint       `json:"business_id"`
	EventType   string     `json:"event_type"`
	Enabled     bool       `json:"enabled"`
	Channels    datatypes.JSON `json:"channels"`
	Filters     datatypes.JSON `json:"filters,omitempty"`
	Description string     `json:"description,omitempty"`
}

// NotificationConfigsListResponse representa la respuesta paginada
type NotificationConfigsListResponse struct {
	Data       []NotificationConfigResponse `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	TotalPages int                          `json:"total_pages"`
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

// IsEventTypeEnabled verifica si un tipo de evento está habilitado para un negocio
func (nc *NotificationConfig) IsEventTypeEnabled(eventType OrderEventType) bool {
	if !nc.Enabled {
		return false
	}
	return nc.EventType == string(eventType)
}

// HasChannel verifica si un canal está habilitado
func (nc *NotificationConfig) HasChannel(channel string) bool {
	// Parsear channels desde JSON
	// Por ahora, asumimos que channels es un array JSON
	// Implementación simplificada - se puede mejorar
	// Si no hay channels configurados, SSE está habilitado por defecto
	return true // Por defecto SSE está habilitado si enabled=true
}


