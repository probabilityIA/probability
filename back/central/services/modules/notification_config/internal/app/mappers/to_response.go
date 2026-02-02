package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ToResponseDTO convierte una entidad a DTO de respuesta
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas + datos relacionados
func ToResponseDTO(entity *entities.IntegrationNotificationConfig) *dtos.NotificationConfigResponseDTO {
	dto := &dtos.NotificationConfigResponseDTO{
		ID:                      entity.ID,
		BusinessID:              entity.BusinessID,
		IntegrationID:           entity.IntegrationID,
		NotificationTypeID:      entity.NotificationTypeID,
		NotificationEventTypeID: entity.NotificationEventTypeID,
		Enabled:                 entity.Enabled,
		Description:             entity.Description,
		OrderStatusIDs:          entity.OrderStatusIDs,
		CreatedAt:               entity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:               entity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Mapear channels desde el campo deprecated (para compatibilidad frontend)
	if len(entity.Channels) > 0 {
		dto.Channels = entity.Channels
	}

	// Mapear event_type desde NotificationEventType.EventCode (nueva forma)
	if entity.NotificationEventType != nil && entity.NotificationEventType.EventCode != "" {
		dto.EventType = &entity.NotificationEventType.EventCode
	} else if entity.EventTypeDeprecated != "" {
		// Fallback a campo deprecated si la relación no está cargada
		dto.EventType = &entity.EventTypeDeprecated
	}

	// Mapear nombres desde relaciones
	if entity.NotificationType != nil {
		dto.NotificationTypeName = &entity.NotificationType.Name
	}

	if entity.NotificationEventType != nil {
		dto.NotificationEventName = &entity.NotificationEventType.EventName
	}

	return dto
}

// ToResponseDTOList convierte una lista de entidades a lista de DTOs
func ToResponseDTOList(entities []entities.IntegrationNotificationConfig) []dtos.NotificationConfigResponseDTO {
	result := make([]dtos.NotificationConfigResponseDTO, 0, len(entities))
	for _, entity := range entities {
		result = append(result, *ToResponseDTO(&entity))
	}
	return result
}
