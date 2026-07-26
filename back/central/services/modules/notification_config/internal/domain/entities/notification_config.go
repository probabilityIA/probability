package entities

import "time"

// IntegrationNotificationConfig representa la configuración de notificaciones para una integración
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas en lugar de strings hardcodeados
type IntegrationNotificationConfig struct {
	ID                      uint
	BusinessID              *uint  // Opcional - auto-asignado por middleware
	IntegrationID           uint   // Integración origen que dispara el evento
	NotificationTypeID      uint   // Canal de salida (WhatsApp, Email, SMS, SSE)
	NotificationEventTypeID uint   // Tipo de evento específico (order.created, order.shipped, etc.)
	Enabled                 bool   // Estado de la configuración
	Description             string // Descripción opcional
	OrderStatusIDs          []uint // IDs de estados de orden (filtro opcional)
	CreatedAt               time.Time
	UpdatedAt               time.Time

	// Campos adicionales para frontend (incluyen datos de relaciones preload)
	EventTypeDeprecated   string              // DEPRECATED - campo event_type de BD (temporal)
	Channels              []string            // DEPRECATED - campo channels JSONB de BD (temporal)
	NotificationType      *NotificationType   // Relación preload con tipo de notificación
	NotificationEventType *NotificationEventType // Relación preload con tipo de evento
}

// ============================================
// LEGACY: Estructuras anteriores (DEPRECADAS)
// Se mantienen temporalmente para compatibilidad
// ============================================

type IntegrationNotificationConfigLegacy struct {
	ID               uint
	IntegrationID    uint
	NotificationType string // "whatsapp", "email", "sms"
	IsActive         bool
	Conditions       NotificationConditions
	Config           NotificationConfig
	Description      string
	Priority         int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NotificationConditions define las condiciones para disparar una notificación (LEGACY)
type NotificationConditions struct {
	Trigger             string   // "order.created", "order.updated", "order.status_changed"
	Statuses            []string // ["pending", "processing"] - vacío = todos
	PaymentMethods      []uint   // [1, 3, 5] - vacío = todos
	SourceIntegrationID *uint    // null = todas las integraciones
}

// NotificationConfig contiene la configuración específica de la notificación (LEGACY)
type NotificationConfig struct {
	TemplateName  string // "confirmacion_pedido_contraentrega"
	RecipientType string // "customer", "business"
	Language      string // "es", "en"
}
