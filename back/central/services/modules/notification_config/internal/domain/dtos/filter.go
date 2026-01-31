package dtos

// FilterNotificationConfigDTO representa los filtros para listar configuraciones
type FilterNotificationConfigDTO struct {
	IntegrationID    *uint
	NotificationType *string
	IsActive         *bool
	Trigger          *string
}
