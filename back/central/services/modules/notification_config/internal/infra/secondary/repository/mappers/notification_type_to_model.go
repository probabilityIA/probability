package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

// NotificationTypeToModel convierte una entidad de dominio a modelo de base de datos
func NotificationTypeToModel(entity *entities.NotificationType) (*models.NotificationType, error) {
	if entity == nil {
		return nil, nil
	}

	var configSchemaJSON datatypes.JSON
	if entity.ConfigSchema != nil {
		jsonBytes, err := json.Marshal(entity.ConfigSchema)
		if err != nil {
			return nil, err
		}
		configSchemaJSON = datatypes.JSON(jsonBytes)
	}

	model := &models.NotificationType{
		Name:         entity.Name,
		Code:         entity.Code,
		Description:  entity.Description,
		Icon:         entity.Icon,
		IsActive:     entity.IsActive,
		ConfigSchema: configSchemaJSON,
	}
	model.ID = entity.ID

	return model, nil
}
