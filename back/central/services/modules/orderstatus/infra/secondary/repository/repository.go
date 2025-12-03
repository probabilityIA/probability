package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// Repository implementa domain.IRepository
type Repository struct {
	db     db.IDatabase
	logger log.ILogger
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, logger log.ILogger) domain.IRepository {
	return &Repository{
		db:     database,
		logger: logger,
	}
}

func (r *Repository) Create(ctx context.Context, mapping *models.OrderStatusMapping) error {
	return r.db.Conn(ctx).Create(mapping).Error
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*models.OrderStatusMapping, error) {
	var mapping models.OrderStatusMapping
	err := r.db.Conn(ctx).Where("id = ?", id).First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order status mapping not found")
		}
		return nil, err
	}
	return &mapping, nil
}

func (r *Repository) List(ctx context.Context, filters map[string]interface{}) ([]models.OrderStatusMapping, int64, error) {
	var mappings []models.OrderStatusMapping
	var total int64

	query := r.db.Conn(ctx).Model(&models.OrderStatusMapping{})

	// Aplicar filtros
	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.Where("integration_type = ?", integrationType)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener resultados
	if err := query.Order("integration_type ASC, priority DESC, created_at DESC").Find(&mappings).Error; err != nil {
		return nil, 0, err
	}

	return mappings, total, nil
}

func (r *Repository) Update(ctx context.Context, mapping *models.OrderStatusMapping) error {
	return r.db.Conn(ctx).Save(mapping).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.OrderStatusMapping{}, id).Error
}

func (r *Repository) ToggleActive(ctx context.Context, id uint) (*models.OrderStatusMapping, error) {
	mapping, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mapping.IsActive = !mapping.IsActive
	if err := r.Update(ctx, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (r *Repository) Exists(ctx context.Context, integrationType, originalStatus string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type = ? AND original_status = ?", integrationType, originalStatus).
		Count(&count).Error
	return count > 0, err
}
