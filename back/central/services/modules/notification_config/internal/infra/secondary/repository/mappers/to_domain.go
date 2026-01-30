package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ToDomain convierte un modelo de base de datos a entidad de dominio
func ToDomain(model *IntegrationNotificationConfigModel) (*entities.IntegrationNotificationConfig, error) {
	if model == nil {
		return nil, nil
	}

	var conditions entities.NotificationConditions
	if err := json.Unmarshal(model.Conditions, &conditions); err != nil {
		return nil, err
	}

	var config entities.NotificationConfig
	if err := json.Unmarshal(model.Config, &config); err != nil {
		return nil, err
	}

	return &entities.IntegrationNotificationConfig{
		ID:               model.ID,
		IntegrationID:    model.IntegrationID,
		NotificationType: model.NotificationType,
		IsActive:         model.IsActive,
		Conditions:       conditions,
		Config:           config,
		Description:      model.Description,
		Priority:         model.Priority,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}, nil
}

// ToDomainList convierte una lista de modelos a lista de entidades
func ToDomainList(models []IntegrationNotificationConfigModel) ([]entities.IntegrationNotificationConfig, error) {
	result := make([]entities.IntegrationNotificationConfig, 0, len(models))

	for _, model := range models {
		entity, err := ToDomain(&model)
		if err != nil {
			return nil, err
		}
		if entity != nil {
			result = append(result, *entity)
		}
	}

	return result, nil
}
