package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type notificationEventTypeRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *notificationEventTypeRepository) GetByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error) {
	var models []models.NotificationEventType

	query := r.db.Conn(ctx).Preload("NotificationType").Preload("AllowedOrderStatuses")
	if notificationTypeID > 0 {
		query = query.Where("notification_type_id = ?", notificationTypeID)
	}

	if err := query.Find(&models).Error; err != nil {
		r.logger.Error().Err(err).Uint("notification_type_id", notificationTypeID).Msg("Error getting notification event types")
		return nil, err
	}

	entities, err := mappers.NotificationEventTypeToDomainList(models)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting models to entities")
		return nil, err
	}

	return entities, nil
}

func (r *notificationEventTypeRepository) GetByID(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
	var model models.NotificationEventType

	if err := r.db.Conn(ctx).Preload("NotificationType").Preload("AllowedOrderStatuses").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotificationEventTypeNotFound
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification event type by ID")
		return nil, err
	}

	entity, err := mappers.NotificationEventTypeToDomain(&model)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting model to entity")
		return nil, err
	}

	return entity, nil
}

func (r *notificationEventTypeRepository) Create(ctx context.Context, eventType *entities.NotificationEventType) error {
	model, err := mappers.NotificationEventTypeToModel(eventType)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error creating notification event type")
		return err
	}

	eventType.ID = model.ID
	if len(eventType.AllowedOrderStatusIDs) > 0 {
		if err := replaceAllowedStatuses(r.db.Conn(ctx), model.ID, eventType.AllowedOrderStatusIDs); err != nil {
			r.logger.Error().Err(err).Uint("id", model.ID).Msg("Error saving allowed order statuses")
			return err
		}
	}
	return nil
}

func (r *notificationEventTypeRepository) Update(ctx context.Context, eventType *entities.NotificationEventType) error {
	model, err := mappers.NotificationEventTypeToModel(eventType)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	result := r.db.Conn(ctx).Model(&models.NotificationEventType{}).
		Where("id = ?", eventType.ID).
		Updates(&model)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Error updating notification event type")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationEventTypeNotFound
	}

	if eventType.AllowedOrderStatusIDs != nil {
		if err := replaceAllowedStatuses(r.db.Conn(ctx), eventType.ID, eventType.AllowedOrderStatusIDs); err != nil {
			r.logger.Error().Err(err).Uint("id", eventType.ID).Msg("Error replacing allowed order statuses")
			return err
		}
	}

	return nil
}

func (r *notificationEventTypeRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Conn(ctx).Unscoped().Delete(&models.NotificationEventType{}, id)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Uint("id", id).Msg("Error deleting notification event type")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationEventTypeNotFound
	}

	return nil
}

func (r *notificationEventTypeRepository) GetAll(ctx context.Context) ([]entities.NotificationEventType, error) {
	r.logger.Info().Msg("🔍 [Repository] Fetching all notification event types from DB")

	var models []models.NotificationEventType

	if err := r.db.Conn(ctx).Preload("NotificationType").Preload("AllowedOrderStatuses").Find(&models).Error; err != nil {
		r.logger.Error().Err(err).Msg("❌ [Repository] Error getting all notification event types from DB")
		return nil, err
	}

	r.logger.Info().Int("count", len(models)).Msg("✅ [Repository] All notification event types fetched from DB")

	entities, err := mappers.NotificationEventTypeToDomainList(models)
	if err != nil {
		r.logger.Error().Err(err).Msg("❌ [Repository] Error converting models to entities")
		return nil, err
	}

	return entities, nil
}

func replaceAllowedStatuses(conn *gorm.DB, eventTypeID uint, statusIDs []uint) error {
	return conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM notification_event_type_allowed_statuses WHERE notification_event_type_id = ?", eventTypeID).Error; err != nil {
			return err
		}
		for _, statusID := range statusIDs {
			if err := tx.Exec(
				"INSERT INTO notification_event_type_allowed_statuses (notification_event_type_id, order_status_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				eventTypeID, statusID,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
