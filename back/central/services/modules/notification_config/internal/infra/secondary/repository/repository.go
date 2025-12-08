package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GormRepository struct {
	db db.IDatabase
}

func New(db db.IDatabase) domain.IRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, config *domain.NotificationConfig) error {
	model, err := toModel(config)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Update ID and timestamps
	config.ID = model.ID
	config.CreatedAt = model.CreatedAt
	config.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*domain.NotificationConfig, error) {
	var model models.BusinessNotificationConfig
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(&model)
}

func (r *GormRepository) Update(ctx context.Context, id uint, config *domain.NotificationConfig) error {
	model, err := toModel(config)
	if err != nil {
		return err
	}
	model.ID = id

	if err := r.db.WithContext(ctx).Model(&models.BusinessNotificationConfig{Model: gorm.Model{ID: id}}).Updates(model).Error; err != nil {
		return err
	}

	// Fetch updated to get timestamps
	updated, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	*config = *updated
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.BusinessNotificationConfig{}, id).Error
}

func (r *GormRepository) List(ctx context.Context, filter domain.ConfigFilter) ([]*domain.NotificationConfig, error) {
	var modelsList []models.BusinessNotificationConfig
	query := r.db.WithContext(ctx)

	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}
	if filter.EventType != nil {
		query = query.Where("event_type = ?", *filter.EventType)
	}

	if err := query.Find(&modelsList).Error; err != nil {
		return nil, err
	}

	configs := make([]*domain.NotificationConfig, len(modelsList))
	for i, m := range modelsList {
		d, err := toDomain(&m)
		if err != nil {
			return nil, err
		}
		configs[i] = d
	}
	return configs, nil
}

func (r *GormRepository) GetByBusinessAndEventType(ctx context.Context, businessID uint, eventType string) (*domain.NotificationConfig, error) {
	var model models.BusinessNotificationConfig
	if err := r.db.WithContext(ctx).Where("business_id = ? AND event_type = ?", businessID, eventType).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(&model)
}

// Helpers

func toModel(d *domain.NotificationConfig) (*models.BusinessNotificationConfig, error) {
	channelsJSON, err := json.Marshal(d.Channels)
	if err != nil {
		return nil, err
	}
	filtersJSON, err := json.Marshal(d.Filters)
	if err != nil {
		return nil, err
	}

	return &models.BusinessNotificationConfig{
		BusinessID:  d.BusinessID,
		EventType:   d.EventType,
		Enabled:     d.Enabled,
		Channels:    datatypes.JSON(channelsJSON),
		Filters:     datatypes.JSON(filtersJSON),
		Description: d.Description,
	}, nil
}

func toDomain(m *models.BusinessNotificationConfig) (*domain.NotificationConfig, error) {
	var channels []string
	if len(m.Channels) > 0 {
		if err := json.Unmarshal(m.Channels, &channels); err != nil {
			return nil, err
		}
	}

	var filters map[string]interface{}
	if len(m.Filters) > 0 {
		if err := json.Unmarshal(m.Filters, &filters); err != nil {
			return nil, err
		}
	}

	deletedAt := m.DeletedAt.Time
	var deletedAtPtr *time.Time
	if m.DeletedAt.Valid {
		deletedAtPtr = &deletedAt
	}

	return &domain.NotificationConfig{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   deletedAtPtr,
		BusinessID:  m.BusinessID,
		EventType:   m.EventType,
		Enabled:     m.Enabled,
		Channels:    channels,
		Filters:     filters,
		Description: m.Description,
	}, nil
}
