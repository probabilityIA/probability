package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

type repository struct {
	db     db.IDatabase
	logger log.ILogger
}

// Create crea una nueva configuración
func (r *repository) Create(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	model, err := mappers.ToModel(config)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error creating notification config")
		return err
	}

	config.ID = model.ID
	config.CreatedAt = model.CreatedAt
	config.UpdatedAt = model.UpdatedAt

	return nil
}

// Update actualiza una configuración existente
func (r *repository) Update(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	model, err := mappers.ToModel(config)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting entity to model")
		return err
	}

	result := r.db.Conn(ctx).Model(&mappers.IntegrationNotificationConfigModel{}).
		Where("id = ?", config.ID).
		Updates(&model)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Error updating notification config")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationConfigNotFound
	}

	return nil
}

// GetByID obtiene una configuración por su ID
func (r *repository) GetByID(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
	var model mappers.IntegrationNotificationConfigModel

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotificationConfigNotFound
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config by ID")
		return nil, err
	}

	entity, err := mappers.ToDomain(&model)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting model to entity")
		return nil, err
	}

	return entity, nil
}

// List obtiene una lista de configuraciones con filtros opcionales
// NUEVA ESTRUCTURA: Filtra por IDs de tablas normalizadas + Preload de relaciones
func (r *repository) List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
	query := r.db.Conn(ctx).Model(&mappers.IntegrationNotificationConfigModel{})

	// PRELOAD de relaciones para traer datos completos
	query = query.Preload("NotificationType").Preload("NotificationEventType")

	// Log de filtros aplicados
	logEvent := r.logger.Info()
	if filters.IntegrationID != nil {
		logEvent = logEvent.Uint("filter_integration_id", *filters.IntegrationID)
		query = query.Where("integration_id = ?", *filters.IntegrationID)
	}
	if filters.NotificationTypeID != nil {
		logEvent = logEvent.Uint("filter_notification_type_id", *filters.NotificationTypeID)
		query = query.Where("notification_type_id = ?", *filters.NotificationTypeID)
	}
	if filters.NotificationEventTypeID != nil {
		logEvent = logEvent.Uint("filter_notification_event_type_id", *filters.NotificationEventTypeID)
		query = query.Where("notification_event_type_id = ?", *filters.NotificationEventTypeID)
	}
	if filters.Enabled != nil {
		logEvent = logEvent.Bool("filter_enabled", *filters.Enabled)
		query = query.Where("enabled = ?", *filters.Enabled)
	}

	// Ordenar por fecha de creación descendente
	query = query.Order("created_at DESC")

	var models []mappers.IntegrationNotificationConfigModel
	if err := query.Find(&models).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing notification configs")
		return nil, err
	}

	entities, err := mappers.ToDomainList(models)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting models to entities")
		return nil, err
	}

	return entities, nil
}

// Delete elimina una configuración por su ID
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.Conn(ctx).Delete(&mappers.IntegrationNotificationConfigModel{}, id)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Uint("id", id).Msg("Error deleting notification config")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrNotificationConfigNotFound
	}

	return nil
}

// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas por integración y trigger
func (r *repository) GetActiveConfigsByIntegrationAndTrigger(
	ctx context.Context,
	integrationID uint,
	trigger string,
) ([]entities.IntegrationNotificationConfig, error) {
	var models []mappers.IntegrationNotificationConfigModel

	err := r.db.Conn(ctx).
		Where("integration_id = ?", integrationID).
		Where("is_active = ?", true).
		Where("conditions->>'trigger' = ?", trigger).
		Order("priority DESC").
		Find(&models).Error

	if err != nil {
		r.logger.Error().
			Err(err).
			Uint("integration_id", integrationID).
			Str("trigger", trigger).
			Msg("Error getting active configs by integration and trigger")
		return nil, err
	}

	entities, err := mappers.ToDomainList(models)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error converting models to entities")
		return nil, err
	}

	return entities, nil
}
