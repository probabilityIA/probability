package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

// NotificationEventTypeToModel convierte una entidad de dominio a modelo de base de datos
func NotificationEventTypeToModel(entity *entities.NotificationEventType) (*models.NotificationEventType, error) {
	if entity == nil {
		return nil, nil
	}

	var templateConfigJSON datatypes.JSON
	if entity.TemplateConfig != nil {
		jsonBytes, err := json.Marshal(entity.TemplateConfig)
		if err != nil {
			return nil, err
		}
		templateConfigJSON = datatypes.JSON(jsonBytes)
	}

	model := &models.NotificationEventType{
		NotificationTypeID: entity.NotificationTypeID,
		EventCode:          entity.EventCode,
		EventName:          entity.EventName,
		Description:        entity.Description,
		TemplateConfig:     templateConfigJSON,
		IsActive:           entity.IsActive,
	}
	model.ID = entity.ID

	return model, nil
}
