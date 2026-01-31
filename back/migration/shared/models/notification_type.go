package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	NOTIFICATION TYPES - Tipos de notificaciones disponibles
//
// ───────────────────────────────────────────

// NotificationType representa un tipo/canal de notificación disponible en el sistema
// Ejemplos: WhatsApp, SSE (Server-Sent Events), Email, SMS
type NotificationType struct {
	gorm.Model

	// Identificación
	Name string `gorm:"size:100;not null"`       // "WhatsApp", "SSE", "Email", "SMS"
	Code string `gorm:"size:50;unique;not null"` // "whatsapp", "sse", "email", "sms"

	// Información
	Description string `gorm:"size:500"` // Descripción del tipo de notificación
	Icon        string `gorm:"size:100"` // Ícono para UI

	// Estado
	IsActive bool `gorm:"default:true;index"` // Si el tipo está activo y disponible

	// Schema de configuración (JSON) - Define qué campos de configuración son necesarios
	// Ejemplo WhatsApp: {"required_fields": ["template_id"], "optional_fields": ["language"]}
	// Ejemplo Email: {"required_fields": ["template"], "optional_fields": ["subject", "reply_to"]}
	ConfigSchema datatypes.JSON `gorm:"type:jsonb"`

	// Relaciones
	NotificationEventTypes []NotificationEventType `gorm:"foreignKey:NotificationTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (NotificationType) TableName() string {
	return "notification_types"
}
