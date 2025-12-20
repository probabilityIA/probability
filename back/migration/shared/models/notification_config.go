package models

import (
	"time"

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

	// Filtros opcionales (JSON) - Para otros tipos de filtros (ej: min_amount, max_amount, etc.)
	// Ya NO se usa para estados de orden (usar OrderStatuses relationship)
	// Ejemplo: {"min_amount": 1000, "currency": "COP"}
	Filters datatypes.JSON `gorm:"type:jsonb"`

	// Descripción opcional
	Description string `gorm:"size:500"`

	// Relación con estados de orden (solo para event_type = "order.status_changed")
	// Si está vacío, se notifican TODOS los estados
	OrderStatuses []OrderStatus `gorm:"many2many:business_notification_config_order_statuses;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (BusinessNotificationConfig) TableName() string {
	return "business_notification_configs"
}

// ───────────────────────────────────────────
//
//	BUSINESS NOTIFICATION CONFIG ORDER STATUS - Tabla intermedia
//
// ───────────────────────────────────────────

// BusinessNotificationConfigOrderStatus representa la relación many-to-many
// entre configuraciones de notificaciones y estados de orden
type BusinessNotificationConfigOrderStatus struct {
	BusinessNotificationConfigID uint `gorm:"primaryKey;not null;index"`
	OrderStatusID                uint `gorm:"primaryKey;not null;index"`

	// Timestamp
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relaciones
	BusinessNotificationConfig BusinessNotificationConfig `gorm:"foreignKey:BusinessNotificationConfigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrderStatus                OrderStatus                `gorm:"foreignKey:OrderStatusID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla intermedia
func (BusinessNotificationConfigOrderStatus) TableName() string {
	return "business_notification_config_order_statuses"
}
