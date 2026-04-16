package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ToDomain convierte un modelo de base de datos a entidad de dominio
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas + relaciones preload
func ToDomain(model *IntegrationNotificationConfigModel) (*entities.IntegrationNotificationConfig, error) {
	if model == nil {
		return nil, nil
	}

	entity := &entities.IntegrationNotificationConfig{
		ID:          model.ID,
		BusinessID:  &model.BusinessID,
		Enabled:     model.Enabled,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	// Mapear IDs (pueden ser nil en modelo)
	if model.IntegrationID != nil {
		entity.IntegrationID = *model.IntegrationID
	}
	if model.NotificationTypeID != nil {
		entity.NotificationTypeID = *model.NotificationTypeID
	}
	if model.NotificationEventTypeID != nil {
		entity.NotificationEventTypeID = *model.NotificationEventTypeID
	}

	// Mapear campos deprecated (temporales)
	entity.EventTypeDeprecated = model.EventType

	// Mapear Channels desde JSONB
	if len(model.Channels) > 0 {
		var channels []string
		if err := json.Unmarshal(model.Channels, &channels); err == nil {
			entity.Channels = channels
		}
	}

	// Mapear relaciones preload (si están cargadas)
	// Nota: En el modelo GORM, las relaciones NO son punteros
	if model.NotificationTypeID != nil && model.NotificationType.ID != 0 {
		entity.NotificationType = &entities.NotificationType{
			ID:   model.NotificationType.ID,
			Name: model.NotificationType.Name,
			Code: model.NotificationType.Code,
		}
	}

	if model.NotificationEventTypeID != nil && model.NotificationEventType.ID != 0 {
		entity.NotificationEventType = &entities.NotificationEventType{
			ID:        model.NotificationEventType.ID,
			EventCode: model.NotificationEventType.EventCode,
			EventName: model.NotificationEventType.EventName,
		}
	}

	// Mapear OrderStatuses M2M a lista de IDs
	if len(model.OrderStatuses) > 0 {
		orderStatusIDs := make([]uint, len(model.OrderStatuses))
		for i, status := range model.OrderStatuses {
			orderStatusIDs[i] = status.ID
		}
		entity.OrderStatusIDs = orderStatusIDs
	}

	return entity, nil
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
