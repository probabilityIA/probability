package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/response"
)

// DomainToResponse convierte una entidad NotificationEventType de dominio a respuesta HTTP
func DomainToResponse(entity entities.NotificationEventType) response.NotificationEventType {
	resp := response.NotificationEventType{
		ID:                 entity.ID,
		NotificationTypeID: entity.NotificationTypeID,
		EventCode:          entity.EventCode,
		EventName:          entity.EventName,
		Description:        entity.Description,
		TemplateConfig:     entity.TemplateConfig,
		IsActive:           entity.IsActive,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
	}

	// Incluir relación si está cargada
	if entity.NotificationType != nil {
		resp.NotificationType = &response.NotificationTypeBasic{
			ID:   entity.NotificationType.ID,
			Name: entity.NotificationType.Name,
			Code: entity.NotificationType.Code,
		}
	}

	return resp
}

// DomainListToResponse convierte una lista de entidades a lista de respuestas HTTP
func DomainListToResponse(entities []entities.NotificationEventType) []response.NotificationEventType {
	results := make([]response.NotificationEventType, len(entities))
	for i, entity := range entities {
		results[i] = DomainToResponse(entity)
	}
	return results
}
