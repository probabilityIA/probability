package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ToResponseDTO convierte una entidad a DTO de respuesta
func ToResponseDTO(entity *entities.IntegrationNotificationConfig) *dtos.NotificationConfigResponseDTO {
	return &dtos.NotificationConfigResponseDTO{
		ID:               entity.ID,
		IntegrationID:    entity.IntegrationID,
		NotificationType: entity.NotificationType,
		IsActive:         entity.IsActive,
		Conditions:       entity.Conditions,
		Config:           entity.Config,
		Description:      entity.Description,
		Priority:         entity.Priority,
		CreatedAt:        entity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        entity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToResponseDTOList convierte una lista de entidades a lista de DTOs
func ToResponseDTOList(entities []entities.IntegrationNotificationConfig) []dtos.NotificationConfigResponseDTO {
	result := make([]dtos.NotificationConfigResponseDTO, 0, len(entities))
	for _, entity := range entities {
		result = append(result, *ToResponseDTO(&entity))
	}
	return result
}
