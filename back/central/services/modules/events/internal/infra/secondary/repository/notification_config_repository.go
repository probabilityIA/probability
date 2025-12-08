package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// NotificationConfigRepository implementa el repositorio de configuraciones de notificaciones
type NotificationConfigRepository struct {
	db db.IDatabase
}

// NewNotificationConfigRepository crea una nueva instancia del repositorio
func New(database db.IDatabase) domain.INotificationConfigRepository {
	return &NotificationConfigRepository{
		db: database,
	}
}

// Create crea una nueva configuración
func (r *NotificationConfigRepository) Create(ctx context.Context, config *domain.NotificationConfig) error {
	dbConfig := mappers.ToDBNotificationConfig(config)
	if err := r.db.Conn(ctx).Create(dbConfig).Error; err != nil {
		return err
	}
	config.ID = dbConfig.ID
	return nil
}

// GetByID obtiene una configuración por su ID
func (r *NotificationConfigRepository) GetByID(ctx context.Context, id uint) (*domain.NotificationConfig, error) {
	var config models.BusinessNotificationConfig
	err := r.db.Conn(ctx).
		Preload("Business").
		Where("id = ?", id).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mappers.ToDomainNotificationConfig(&config), nil
}

// GetByBusinessAndEventType obtiene una configuración por business_id y event_type
func (r *NotificationConfigRepository) GetByBusinessAndEventType(ctx context.Context, businessID uint, eventType string) (*domain.NotificationConfig, error) {
	var config models.BusinessNotificationConfig
	err := r.db.Conn(ctx).
		Preload("Business").
		Where("business_id = ? AND event_type = ?", businessID, eventType).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No existe configuración, retornar nil
		}
		return nil, err
	}

	return mappers.ToDomainNotificationConfig(&config), nil
}

// GetByBusinessID obtiene todas las configuraciones de un negocio
func (r *NotificationConfigRepository) GetByBusinessID(ctx context.Context, businessID uint) ([]domain.NotificationConfig, error) {
	var configs []models.BusinessNotificationConfig
	err := r.db.Conn(ctx).
		Preload("Business").
		Where("business_id = ?", businessID).
		Order("event_type ASC").
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	domainConfigs := make([]domain.NotificationConfig, len(configs))
	for i, config := range configs {
		domainConfigs[i] = *mappers.ToDomainNotificationConfig(&config)
	}

	return domainConfigs, nil
}

// List obtiene una lista paginada de configuraciones
func (r *NotificationConfigRepository) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.NotificationConfig, int64, error) {
	var configs []models.BusinessNotificationConfig
	var total int64

	query := r.db.Conn(ctx).Model(&models.BusinessNotificationConfig{})

	// Filtros
	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.Where("business_id = ?", businessID)
	}

	if eventType, ok := filters["event_type"].(string); ok && eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	if enabled, ok := filters["enabled"].(bool); ok {
		query = query.Where("enabled = ?", enabled)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Precargar relaciones
	query = query.Preload("Business")

	// Ejecutar query
	if err := query.Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a dominio
	domainConfigs := make([]domain.NotificationConfig, len(configs))
	for i, config := range configs {
		domainConfigs[i] = *mappers.ToDomainNotificationConfig(&config)
	}

	return domainConfigs, total, nil
}

// Update actualiza una configuración
func (r *NotificationConfigRepository) Update(ctx context.Context, config *domain.NotificationConfig) error {
	dbConfig := mappers.ToDBNotificationConfig(config)
	return r.db.Conn(ctx).Save(dbConfig).Error
}

// Delete elimina una configuración
func (r *NotificationConfigRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Where("id = ?", id).Delete(&models.BusinessNotificationConfig{}).Error
}

// IsEventTypeEnabled verifica si un tipo de evento está habilitado para un negocio
func (r *NotificationConfigRepository) IsEventTypeEnabled(ctx context.Context, businessID uint, eventType string) (bool, error) {
	config, err := r.GetByBusinessAndEventType(ctx, businessID, eventType)
	if err != nil {
		return false, err
	}

	// Si no existe configuración, retornar true por defecto
	if config == nil {
		return true, nil
	}

	return config.Enabled, nil
}
