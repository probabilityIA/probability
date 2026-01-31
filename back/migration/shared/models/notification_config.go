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
	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// NUEVA ARQUITECTURA: Relación con Integration (la integración que genera el evento)
	// Esta es la integración de ORIGEN que dispara la notificación
	// Ejemplo: Integration "Shopify - Mi Tiendita" (id: 5) genera evento order.created
	// NULLABLE durante migración, se populará con datos existentes
	IntegrationID *uint       `gorm:"index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// NUEVA ARQUITECTURA: Relación con NotificationType (canal de salida)
	// Define POR DÓNDE se envía la notificación (WhatsApp, SSE, Email, SMS)
	// NULLABLE durante migración, se populará con datos existentes
	NotificationTypeID *uint            `gorm:"index"`
	NotificationType   NotificationType `gorm:"foreignKey:NotificationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// NUEVA ARQUITECTURA: Relación con NotificationEventType (tipo de evento específico)
	// Define QUÉ EVENTO dispara la notificación (order.created, order.shipped, etc.)
	// NULLABLE durante migración, se populará con datos existentes
	NotificationEventTypeID *uint                 `gorm:"index"`
	NotificationEventType   NotificationEventType `gorm:"foreignKey:NotificationEventTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Si la notificación está habilitada
	Enabled bool `gorm:"default:true;index"`

	// DEPRECATED: EventType - Ahora se usa NotificationEventType relationship
	// Mantener temporalmente para migración de datos existentes
	EventType string `gorm:"size:64;index"`

	// DEPRECATED: Channels - Ahora se usa NotificationType relationship
	// Esta columna será eliminada después de la migración
	// Channels datatypes.JSON `gorm:"type:jsonb"` // REMOVIDO

	// Filtros opcionales (JSON) - Para otros tipos de filtros (ej: min_amount, max_amount, etc.)
	// Ya NO se usa para estados de orden (usar OrderStatuses relationship)
	// Ejemplo: {"min_amount": 1000, "currency": "COP"}
	Filters datatypes.JSON `gorm:"type:jsonb"`

	// Descripción opcional
	Description string `gorm:"size:500"`

	// Relación con estados de orden (M2M)
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
