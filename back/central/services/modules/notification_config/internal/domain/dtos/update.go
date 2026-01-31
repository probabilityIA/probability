package dtos

import "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"

// UpdateNotificationConfigDTO representa los datos para actualizar una configuración
type UpdateNotificationConfigDTO struct {
	NotificationType *string
	IsActive         *bool
	Conditions       *entities.NotificationConditions
	Config           *entities.NotificationConfig
	Description      *string
	Priority         *int
}
