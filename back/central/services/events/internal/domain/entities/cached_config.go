package entities

// CachedNotificationConfig representa una configuración de notificación leída del cache Redis.
// Espejo de la estructura cacheada por el módulo notification_config.
type CachedNotificationConfig struct {
	ID                      uint
	BusinessID              *uint
	IntegrationID           uint
	NotificationTypeID      uint   // 1=SSE, 2=WhatsApp, 3=Email
	NotificationEventTypeID uint
	Enabled                 bool
	Description             string
	OrderStatusIDs          []uint
	EventCode               string
	OrderStatusCodes        []string
}
