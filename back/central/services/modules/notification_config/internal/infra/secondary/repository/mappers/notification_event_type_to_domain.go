package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// NotificationEventTypeToDomain convierte un modelo de base de datos a entidad de dominio
func NotificationEventTypeToDomain(model *models.NotificationEventType) (*entities.NotificationEventType, error) {
	if model == nil {
		return nil, nil
	}

	var templateConfig map[string]interface{}
	if model.TemplateConfig != nil {
		if err := json.Unmarshal(model.TemplateConfig, &templateConfig); err != nil {
			return nil, err
		}
	}

	entity := &entities.NotificationEventType{
		ID:                 model.ID,
		NotificationTypeID: model.NotificationTypeID,
		EventCode:          model.EventCode,
		EventName:          model.EventName,
		Description:        model.Description,
		TemplateConfig:     templateConfig,
		IsActive:           model.IsActive,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}

	// Convertir relación NotificationType si está pre-cargada
	if model.NotificationType.ID != 0 {
		notificationType, err := NotificationTypeToDomain(&model.NotificationType)
		if err != nil {
			return nil, err
		}
		entity.NotificationType = notificationType
	}

	return entity, nil
}

// NotificationEventTypeToDomainList convierte una lista de modelos a lista de entidades
func NotificationEventTypeToDomainList(modelsList []models.NotificationEventType) ([]entities.NotificationEventType, error) {
	result := make([]entities.NotificationEventType, 0, len(modelsList))

	for _, model := range modelsList {
		entity, err := NotificationEventTypeToDomain(&model)
		if err != nil {
			return nil, err
		}
		if entity != nil {
			result = append(result, *entity)
		}
	}

	return result, nil
}
