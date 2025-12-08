package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBNotificationConfig convierte una configuración de dominio a modelo de base de datos
func ToDBNotificationConfig(nc *domain.NotificationConfig) *models.BusinessNotificationConfig {
	if nc == nil {
		return nil
	}
	dbConfig := &models.BusinessNotificationConfig{
		Model: gorm.Model{
			ID:        nc.ID,
			CreatedAt: nc.CreatedAt,
			UpdatedAt: nc.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		BusinessID:  nc.BusinessID,
		EventType:   nc.EventType,
		Enabled:     nc.Enabled,
		Channels:    nc.Channels,
		Filters:     nc.Filters,
		Description: nc.Description,
	}
	if nc.DeletedAt != nil {
		dbConfig.DeletedAt = gorm.DeletedAt{Time: *nc.DeletedAt, Valid: true}
	}
	return dbConfig
}

// ToDomainNotificationConfig convierte una configuración de base de datos a dominio
func ToDomainNotificationConfig(nc *models.BusinessNotificationConfig) *domain.NotificationConfig {
	if nc == nil {
		return nil
	}
	var deletedAt *time.Time
	if nc.DeletedAt.Valid {
		deletedAt = &nc.DeletedAt.Time
	}
	return &domain.NotificationConfig{
		ID:          nc.ID,
		CreatedAt:   nc.CreatedAt,
		UpdatedAt:   nc.UpdatedAt,
		DeletedAt:   deletedAt,
		BusinessID:  nc.BusinessID,
		EventType:   nc.EventType,
		Enabled:     nc.Enabled,
		Channels:    nc.Channels,
		Filters:     nc.Filters,
		Description: nc.Description,
	}
}

