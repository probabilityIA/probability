package dtos

// FilterNotificationConfigDTO representa los filtros para listar configuraciones
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
type FilterNotificationConfigDTO struct {
	IntegrationID           *uint
	NotificationTypeID      *uint
	NotificationEventTypeID *uint
	Enabled                 *bool
}

// ============================================
// LEGACY: Filtros anteriores (DEPRECADOS)
// Se mantienen temporalmente para compatibilidad
// ============================================

type FilterNotificationConfigDTOLegacy struct {
	IntegrationID    *uint
	NotificationType *string
	IsActive         *bool
	Trigger          *string
}
