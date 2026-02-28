package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// SyncConfigs ejecuta create/update/delete en una transacción atómica
func (r *repository) SyncConfigs(
	ctx context.Context,
	businessID uint,
	integrationID uint,
	toCreate []*entities.IntegrationNotificationConfig,
	toUpdate []*entities.IntegrationNotificationConfig,
	toDeleteIDs []uint,
) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. DELETE: soft-delete configs que ya no están en el request
		for _, id := range toDeleteIDs {
			// Limpiar M2M antes de soft-delete
			if err := tx.Model(&mappers.IntegrationNotificationConfigModel{Model: gorm.Model{ID: id}}).
				Association("OrderStatuses").Clear(); err != nil {
				r.logger.Error().Err(err).Uint("id", id).Msg("Error clearing order statuses before delete")
				return fmt.Errorf("failed to clear order statuses for config %d: %w", id, err)
			}
			if err := tx.Delete(&mappers.IntegrationNotificationConfigModel{}, id).Error; err != nil {
				r.logger.Error().Err(err).Uint("id", id).Msg("Error deleting config in sync")
				return fmt.Errorf("failed to delete config %d: %w", id, err)
			}
		}

		// 2. CREATE: insertar nuevas configs
		for _, entity := range toCreate {
			model, err := mappers.ToModel(entity)
			if err != nil {
				return fmt.Errorf("failed to convert entity to model for create: %w", err)
			}

			if err := tx.Create(model).Error; err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}

			// Reemplazar M2M de order statuses
			if len(entity.OrderStatusIDs) > 0 {
				orderStatuses := make([]models.OrderStatus, len(entity.OrderStatusIDs))
				for i, sid := range entity.OrderStatusIDs {
					orderStatuses[i].ID = sid
				}
				if err := tx.Model(model).Association("OrderStatuses").Replace(orderStatuses); err != nil {
					return fmt.Errorf("failed to set order statuses for new config: %w", err)
				}
			}

			entity.ID = model.ID
			entity.CreatedAt = model.CreatedAt
			entity.UpdatedAt = model.UpdatedAt
		}

		// 3. UPDATE: actualizar configs existentes
		for _, entity := range toUpdate {
			model, err := mappers.ToModel(entity)
			if err != nil {
				return fmt.Errorf("failed to convert entity to model for update: %w", err)
			}

			// Actualizar campos
			if err := tx.Model(&mappers.IntegrationNotificationConfigModel{}).
				Where("id = ?", entity.ID).
				Updates(map[string]interface{}{
					"notification_type_id":       model.NotificationTypeID,
					"notification_event_type_id": model.NotificationEventTypeID,
					"enabled":                    model.Enabled,
					"description":                model.Description,
				}).Error; err != nil {
				return fmt.Errorf("failed to update config %d: %w", entity.ID, err)
			}

			// Reemplazar M2M de order statuses
			configModel := &mappers.IntegrationNotificationConfigModel{}
			configModel.ID = entity.ID
			orderStatuses := make([]models.OrderStatus, len(entity.OrderStatusIDs))
			for i, sid := range entity.OrderStatusIDs {
				orderStatuses[i].ID = sid
			}
			if err := tx.Model(configModel).Association("OrderStatuses").Replace(orderStatuses); err != nil {
				return fmt.Errorf("failed to replace order statuses for config %d: %w", entity.ID, err)
			}
		}

		return nil
	})
}
