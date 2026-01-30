package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"gorm.io/datatypes"
)

// ToModel convierte una entidad de dominio a modelo de base de datos
func ToModel(entity *entities.IntegrationNotificationConfig) (*IntegrationNotificationConfigModel, error) {
	if entity == nil {
		return nil, nil
	}

	conditionsJSON, err := json.Marshal(entity.Conditions)
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(entity.Config)
	if err != nil {
		return nil, err
	}

	return &IntegrationNotificationConfigModel{
		ID:               entity.ID,
		IntegrationID:    entity.IntegrationID,
		NotificationType: entity.NotificationType,
		IsActive:         entity.IsActive,
		Conditions:       datatypes.JSON(conditionsJSON),
		Config:           datatypes.JSON(configJSON),
		Description:      entity.Description,
		Priority:         entity.Priority,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}, nil
}
