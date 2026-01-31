package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	NOTIFICATION EVENT TYPES - Tipos de eventos de notificación
//
// ───────────────────────────────────────────

// NotificationEventType representa un tipo de evento específico para un tipo de notificación
// Ejemplos: "order.created", "order.shipped", "invoice.created"
// Cada NotificationEventType pertenece a un NotificationType (ej: WhatsApp, Email)
type NotificationEventType struct {
	gorm.Model

	// Relación con NotificationType
	NotificationTypeID uint             `gorm:"not null;index;uniqueIndex:idx_notification_type_event_code,priority:1"`
	NotificationType   NotificationType `gorm:"foreignKey:NotificationTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Identificación del evento
	EventCode string `gorm:"size:100;not null;index;uniqueIndex:idx_notification_type_event_code,priority:2"` // "order.created", "order.shipped"
	EventName string `gorm:"size:200;not null"`                                                                 // "Confirmación de Pedido", "Pedido Enviado"

	// Información
	Description string `gorm:"type:text"` // Descripción del evento

	// Configuración de template (JSON) - Define configuraciones específicas del template
	// Ejemplo WhatsApp: {"default_template_id": "order_confirmation_v2", "template_language": "es"}
	// Ejemplo Email: {"default_subject": "Tu pedido ha sido confirmado", "template_file": "order_confirmation.html"}
	TemplateConfig datatypes.JSON `gorm:"type:jsonb"`

	// Estado
	IsActive bool `gorm:"default:true;index"` // Si el evento está activo

	// Relaciones
	BusinessNotificationConfigs []BusinessNotificationConfig `gorm:"foreignKey:NotificationEventTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (NotificationEventType) TableName() string {
	return "notification_event_types"
}
