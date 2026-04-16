package entities

import "time"

// BusinessNotificationConfig representa la configuración de notificaciones para un negocio
// Entidad PURA del dominio - SIN tags de frameworks
// Definela QUÉ notificaciones se envían, POR QUÉ CANAL y EN QUÉ ESTADOS
type BusinessNotificationConfig struct {
	ID                      uint
	BusinessID              uint
	IntegrationID           uint // La integración que genera el evento (ej: Shopify Store 1)
	NotificationTypeID      uint // Canal de salida (WhatsApp, SSE, Email, SMS)
	NotificationEventTypeID uint // Tipo de evento (order.created, order.shipped, etc.)
	Enabled                 bool
	Filters                 map[string]interface{} // Filtros adicionales (ej: min_amount, currency)
	Description             string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time

	// Relaciones (pueden ser nil si no se cargan)
	NotificationType      *NotificationType
	NotificationEventType *NotificationEventType
	OrderStatusIDs        []uint // IDs de OrderStatus a notificar (relación M2M)

	// Campos deprecados (mantener para migración, pero NO usar en nueva lógica)
	EventTypeDeprecated string   // DEPRECATED - ahora se usa NotificationEventType
	Channels            []string // DEPRECATED - ahora se usa NotificationType (temporal para frontend)
}

// OrderStatus representa un estado de orden (entidad relacionada)
// Entidad PURA del dominio - SIN tags de frameworks
type OrderStatus struct {
	ID          uint
	Code        string // "pending", "processing", "completed", etc.
	Name        string // "Pendiente", "En Procesamiento", "Completada"
	Description string
	Category    string // "active", "completed", "cancelled", "refunded"
	IsActive    bool
	Icon        string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
