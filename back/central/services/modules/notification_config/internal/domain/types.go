package domain

import (
	"time"
)

// NotificationConfig representa la configuración de notificaciones para un negocio
type NotificationConfig struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	BusinessID  uint                   `json:"business_id"`
	EventType   string                 `json:"event_type"`
	Enabled     bool                   `json:"enabled"`
	Channels    []string               `json:"channels"`
	Filters     map[string]interface{} `json:"filters"`
	Description string                 `json:"description"`
}

// CreateConfigDTO datos para crear una configuración
type CreateConfigDTO struct {
	BusinessID  uint                   `json:"business_id" binding:"required"`
	EventType   string                 `json:"event_type" binding:"required"`
	Enabled     bool                   `json:"enabled"`
	Channels    []string               `json:"channels"`
	Filters     map[string]interface{} `json:"filters"`
	Description string                 `json:"description"`
}

// UpdateConfigDTO datos para actualizar una configuración
type UpdateConfigDTO struct {
	Enabled     *bool                  `json:"enabled,omitempty"`
	Channels    []string               `json:"channels,omitempty"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	Description *string                `json:"description,omitempty"`
}

// ConfigFilter filtros para listar configuraciones
type ConfigFilter struct {
	BusinessID *uint   `json:"business_id,omitempty" form:"business_id"`
	EventType  *string `json:"event_type,omitempty" form:"event_type"`
}
