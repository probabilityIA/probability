package entities

import "time"

// IntegrationNotificationConfig representa la configuración de notificaciones para una integración
type IntegrationNotificationConfig struct {
	ID               uint
	IntegrationID    uint
	NotificationType string // "whatsapp", "email", "sms"
	IsActive         bool
	Conditions       NotificationConditions
	Config           NotificationConfig
	Description      string
	Priority         int // Mayor prioridad se evalúa primero
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NotificationConditions define las condiciones para disparar una notificación
type NotificationConditions struct {
	Trigger             string   // "order.created", "order.updated", "order.status_changed"
	Statuses            []string // ["pending", "processing"] - vacío = todos
	PaymentMethods      []uint   // [1, 3, 5] - vacío = todos
	SourceIntegrationID *uint    // null = todas las integraciones
}

// NotificationConfig contiene la configuración específica de la notificación
type NotificationConfig struct {
	TemplateName  string // "confirmacion_pedido_contraentrega"
	RecipientType string // "customer", "business"
	Language      string // "es", "en"
}
