package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	BUSINESS NOTIFICATION CONFIG - Configuraciones de notificaciones internas por negocio
//
// ───────────────────────────────────────────

// BusinessNotificationConfig configura qué eventos de órdenes se notifican a un negocio
// Estas son notificaciones internas para el panel administrativo (SSE)
type BusinessNotificationConfig struct {
	gorm.Model

	// Relación con Business
	BusinessID uint     `gorm:"not null;index;uniqueIndex:idx_business_event_type,priority:1"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Tipo de evento de orden que se notifica
	// "order.created", "order.status_changed", "order.cancelled", etc.
	EventType string `gorm:"size:64;not null;index;uniqueIndex:idx_business_event_type,priority:2"`

	// Si la notificación está habilitada
	Enabled bool `gorm:"default:true;index"`

	// Canales de notificación habilitados (JSON array)
	// ["sse", "email", "webhook"] - por ahora solo SSE para notificaciones internas
	Channels datatypes.JSON `gorm:"type:jsonb"`

	// Filtros opcionales (JSON)
	// Permite filtrar eventos por condiciones específicas
	// Ejemplo: {"statuses": ["pending", "processing"], "min_amount": 1000}
	Filters datatypes.JSON `gorm:"type:jsonb"`

	// Descripción opcional
	Description string `gorm:"size:500"`
}

// TableName especifica el nombre de la tabla
func (BusinessNotificationConfig) TableName() string {
	return "business_notification_configs"
}

