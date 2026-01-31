package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// NotificationTypeToDomain convierte un modelo de base de datos a entidad de dominio
func NotificationTypeToDomain(model *models.NotificationType) (*entities.NotificationType, error) {
	if model == nil {
		return nil, nil
	}

	var configSchema map[string]interface{}
	if model.ConfigSchema != nil {
		if err := json.Unmarshal(model.ConfigSchema, &configSchema); err != nil {
			return nil, err
		}
	}

	return &entities.NotificationType{
		ID:           model.ID,
		Name:         model.Name,
		Code:         model.Code,
		Description:  model.Description,
		Icon:         model.Icon,
		IsActive:     model.IsActive,
		ConfigSchema: configSchema,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

// NotificationTypeToDomainList convierte una lista de modelos a lista de entidades
func NotificationTypeToDomainList(modelsList []models.NotificationType) ([]entities.NotificationType, error) {
	result := make([]entities.NotificationType, 0, len(modelsList))

	for _, model := range modelsList {
		entity, err := NotificationTypeToDomain(&model)
		if err != nil {
			return nil, err
		}
		if entity != nil {
			result = append(result, *entity)
		}
	}

	return result, nil
}
