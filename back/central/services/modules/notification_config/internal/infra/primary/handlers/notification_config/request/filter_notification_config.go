package request

// FilterNotificationConfig representa los query params HTTP para filtros
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
type FilterNotificationConfig struct {
	IntegrationID           *uint `form:"integration_id"`
	NotificationTypeID      *uint `form:"notification_type_id"`
	NotificationEventTypeID *uint `form:"notification_event_type_id"`
	Enabled                 *bool `form:"enabled"`
}
