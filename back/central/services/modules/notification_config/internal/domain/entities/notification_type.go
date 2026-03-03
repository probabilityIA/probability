package entities

import "time"

// NotificationType representa un tipo/canal de notificación disponible en el sistema
// Entidad PURA del dominio - SIN tags de frameworks
// Ejemplos: WhatsApp, SSE, Email, SMS
type NotificationType struct {
	ID           uint
	Name         string // "WhatsApp", "SSE", "Email", "SMS"
	Code         string // "whatsapp", "sse", "email", "sms"
	Description  string
	Icon         string
	IsActive     bool
	ConfigSchema map[string]interface{} // Schema de configuración requerida
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NotificationEventType representa un tipo de evento específico para un tipo de notificación
// Entidad PURA del dominio - SIN tags de frameworks
// Ejemplos: "order.created", "order.shipped", "invoice.created"
type NotificationEventType struct {
	ID                 uint
	NotificationTypeID uint
	EventCode          string // "order.created", "order.shipped"
	EventName          string // "Confirmación de Pedido", "Pedido Enviado"
	Description        string
	TemplateConfig     map[string]interface{} // Configuración de template
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relación (opcional, puede ser nil)
	NotificationType *NotificationType

	// Estados de orden permitidos para este tipo de evento
	// Si está vacío → todos los estados están permitidos
	AllowedOrderStatusIDs []uint
}
