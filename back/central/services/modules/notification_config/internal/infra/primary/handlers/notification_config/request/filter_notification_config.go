package request

// FilterNotificationConfig representa los query params HTTP para filtros
type FilterNotificationConfig struct {
	IntegrationID    *uint   `form:"integration_id"`
	NotificationType *string `form:"notification_type"`
	IsActive         *bool   `form:"is_active"`
	Trigger          *string `form:"trigger"`
}
