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

// GetByNotificationType obtiene todos los eventos de un tipo de notificación
func (r *notificationEventTypeRepository) GetByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error) {
	var models []models.NotificationEventType

	query := r.db.Conn(ctx).Preload("NotificationType")
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

// GetByID obtiene un evento de notificación por su ID
func (r *notificationEventTypeRepository) GetByID(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
	var model models.NotificationEventType

	if err := r.db.Conn(ctx).Preload("NotificationType").First(&model, id).Error; err != nil {
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

// Create crea un nuevo evento de notificación
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
	return nil
}

// Update actualiza un evento de notificación existente
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

	return nil
}

// Delete elimina un evento de notificación por su ID (soft delete)
func (r *notificationEventTypeRepository) Delete(ctx context.Context, id uint) error {
	// Unscoped() hace que la eliminación sea permanente (hard delete) en lugar de soft delete
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

// GetAll obtiene todos los eventos de notificación sin filtros
func (r *notificationEventTypeRepository) GetAll(ctx context.Context) ([]entities.NotificationEventType, error) {
	r.logger.Info().Msg("🔍 [Repository] Fetching all notification event types from DB")

	var models []models.NotificationEventType

	if err := r.db.Conn(ctx).Preload("NotificationType").Find(&models).Error; err != nil {
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
