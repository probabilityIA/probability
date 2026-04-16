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

type notificationTypeRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

// GetAll obtiene todos los tipos de notificaciones
func (r *notificationTypeRepository) GetAll(ctx context.Context) ([]entities.NotificationType, error) {
	var models []models.NotificationType

	if err := r.db.Conn(ctx).Find(&models).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error getting all notification types")
		return nil, err
	}

	entities, err := mappers.NotificationTypeToDomainList(models)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting models to entities")
		return nil, err
	}

	return entities, nil
}

// GetByID obtiene un tipo de notificación por su ID
func (r *notificationTypeRepository) GetByID(ctx context.Context, id uint) (*entities.NotificationType, error) {
	var model models.NotificationType

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotificationTypeNotFound
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification type by ID")
		return nil, err
	}

	entity, err := mappers.NotificationTypeToDomain(&model)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting model to entity")
		return nil, err
	}

	return entity, nil
}

// GetByCode obtiene un tipo de notificación por su código
func (r *notificationTypeRepository) GetByCode(ctx context.Context, code string) (*entities.NotificationType, error) {
	var model models.NotificationType

	if err := r.db.Conn(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotificationTypeNotFound
		}
		r.logger.Error().Err(err).Str("code", code).Msg("Error getting notification type by code")
		return nil, err
	}

	entity, err := mappers.NotificationTypeToDomain(&model)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting model to entity")
		return nil, err
	}

	return entity, nil
}

// Create crea un nuevo tipo de notificación
func (r *notificationTypeRepository) Create(ctx context.Context, notificationType *entities.NotificationType) error {
	model, err := mappers.NotificationTypeToModel(notificationType)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error creating notification type")
		return err
	}

	notificationType.ID = model.ID
	notificationType.CreatedAt = model.CreatedAt
	notificationType.UpdatedAt = model.UpdatedAt

	return nil
}

// Update actualiza un tipo de notificación existente
func (r *notificationTypeRepository) Update(ctx context.Context, notificationType *entities.NotificationType) error {
	model, err := mappers.NotificationTypeToModel(notificationType)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	result := r.db.Conn(ctx).Model(&models.NotificationType{}).
		Where("id = ?", notificationType.ID).
		Updates(&model)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Error updating notification type")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationTypeNotFound
	}

	return nil
}

// Delete elimina un tipo de notificación por su ID (soft delete)
func (r *notificationTypeRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Conn(ctx).Delete(&models.NotificationType{}, id)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Uint("id", id).Msg("Error deleting notification type")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationTypeNotFound
	}

	return nil
}
