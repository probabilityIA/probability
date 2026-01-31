package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/response"
)

// DomainToResponse convierte una entidad NotificationType de dominio a respuesta HTTP
func DomainToResponse(entity entities.NotificationType) response.NotificationType {
	return response.NotificationType{
		ID:           entity.ID,
		Name:         entity.Name,
		Code:         entity.Code,
		Description:  entity.Description,
		Icon:         entity.Icon,
		IsActive:     entity.IsActive,
		ConfigSchema: entity.ConfigSchema,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

// DomainListToResponse convierte una lista de entidades a lista de respuestas HTTP
func DomainListToResponse(entities []entities.NotificationType) []response.NotificationType {
	results := make([]response.NotificationType, len(entities))
	for i, entity := range entities {
		results[i] = DomainToResponse(entity)
	}
	return results
}
