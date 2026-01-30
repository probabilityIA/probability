package dtos

import "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"

// CreateNotificationConfigDTO representa los datos para crear una configuración
type CreateNotificationConfigDTO struct {
	IntegrationID    uint
	NotificationType string
	IsActive         bool
	Conditions       entities.NotificationConditions
	Config           entities.NotificationConfig
	Description      string
	Priority         int
}
