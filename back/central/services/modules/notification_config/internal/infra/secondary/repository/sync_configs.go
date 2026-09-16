package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository/mappers"
	"gorm.io/gorm"
)

func (r *repository) SyncConfigs(
	ctx context.Context,
	businessID uint,
	integrationID uint,
	toCreate []*entities.IntegrationNotificationConfig,
	toUpdate []*entities.IntegrationNotificationConfig,
	toDeleteIDs []uint,
) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range toDeleteIDs {
			if err := replaceConfigStatuses(tx, id, nil); err != nil {
				r.logger.Error().Err(err).Uint("id", id).Msg("Error clearing order statuses before delete")
				return fmt.Errorf("failed to clear order statuses for config %d: %w", id, err)
			}
			if err := tx.Delete(&mappers.IntegrationNotificationConfigModel{}, id).Error; err != nil {
				r.logger.Error().Err(err).Uint("id", id).Msg("Error deleting config in sync")
				return fmt.Errorf("failed to delete config %d: %w", id, err)
			}
		}

		for _, entity := range toCreate {
			model, err := mappers.ToModel(entity)
			if err != nil {
				return fmt.Errorf("failed to convert entity to model for create: %w", err)
			}

			if err := tx.Create(model).Error; err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}

			if err := replaceConfigStatuses(tx, model.ID, entity.OrderStatusIDs); err != nil {
				return fmt.Errorf("failed to set order statuses for new config: %w", err)
			}

			entity.ID = model.ID
			entity.CreatedAt = model.CreatedAt
			entity.UpdatedAt = model.UpdatedAt
		}

		for _, entity := range toUpdate {
			model, err := mappers.ToModel(entity)
			if err != nil {
				return fmt.Errorf("failed to convert entity to model for update: %w", err)
			}

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

			if err := replaceConfigStatuses(tx, entity.ID, entity.OrderStatusIDs); err != nil {
				return fmt.Errorf("failed to replace order statuses for config %d: %w", entity.ID, err)
			}
		}

		return nil
	})
}

func replaceConfigStatuses(tx *gorm.DB, configID uint, statusIDs []uint) error {
	if err := tx.Exec("DELETE FROM business_notification_config_order_statuses WHERE business_notification_config_id = ?", configID).Error; err != nil {
		return err
	}
	for _, statusID := range statusIDs {
		if err := tx.Exec(
			"INSERT INTO business_notification_config_order_statuses (business_notification_config_id, order_status_id, created_at) VALUES (?, ?, NOW()) ON CONFLICT DO NOTHING",
			configID, statusID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}
