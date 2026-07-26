package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ToModel convierte una entidad IntegrationNotificationConfig a modelo de BD
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
func ToModel(entity *entities.IntegrationNotificationConfig) (*IntegrationNotificationConfigModel, error) {
	if entity == nil {
		return nil, nil
	}

	model := &IntegrationNotificationConfigModel{
		IntegrationID:           &entity.IntegrationID,
		NotificationTypeID:      &entity.NotificationTypeID,
		NotificationEventTypeID: &entity.NotificationEventTypeID,
		Enabled:                 entity.Enabled,
		Description:             entity.Description,
	}

	// Asignar BusinessID si no es nil
	if entity.BusinessID != nil {
		model.BusinessID = *entity.BusinessID
	}

	// Asignar ID si existe (para updates)
	if entity.ID > 0 {
		model.ID = entity.ID
	}

	// Mapear OrderStatusIDs a relación M2M (si se proporcionan)
	// Solo asignar si hay IDs, para evitar sobrescribir relaciones existentes
	if len(entity.OrderStatusIDs) > 0 {
		orderStatuses := make([]models.OrderStatus, len(entity.OrderStatusIDs))
		for i, statusID := range entity.OrderStatusIDs {
			// Solo asignar el ID, GORM se encargará de la relación
			orderStatuses[i].ID = statusID
		}
		model.OrderStatuses = orderStatuses
	}

	return model, nil
}
