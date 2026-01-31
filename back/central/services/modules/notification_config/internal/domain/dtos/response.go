package dtos

import "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"

// NotificationConfigResponseDTO representa la respuesta de una configuración
type NotificationConfigResponseDTO struct {
	ID               uint
	IntegrationID    uint
	NotificationType string
	IsActive         bool
	Conditions       entities.NotificationConditions
	Config           entities.NotificationConfig
	Description      string
	Priority         int
	CreatedAt        string
	UpdatedAt        string
}
