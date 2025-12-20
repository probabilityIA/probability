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
	err := r.db.Conn(ctx).Preload("IntegrationType").Preload("OrderStatus").Where("id = ?", id).First(&mapping).Error
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

	query := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).Preload("IntegrationType").Preload("OrderStatus")

	// Aplicar filtros
	if integrationTypeID, ok := filters["integration_type_id"].(uint); ok && integrationTypeID > 0 {
		query = query.Where("integration_type_id = ?", integrationTypeID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Contar total antes de aplicar paginación
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	if page, ok := filters["page"].(int); ok && page > 0 {
		if pageSize, ok := filters["page_size"].(int); ok && pageSize > 0 {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}
	}

	// Obtener resultados
	if err := query.Order("integration_type_id ASC, priority DESC, created_at DESC").Find(&mappings).Error; err != nil {
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

func (r *Repository) Exists(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type_id = ? AND original_status = ? AND is_active = ?", integrationTypeID, originalStatus, true).
		Count(&count).Error
	return count > 0, err
}

// GetOrderStatusIDByIntegrationTypeAndOriginalStatus obtiene el order_status_id mapeado para un integration_type_id y original_status dado
// Retorna nil si no se encuentra el mapeo
// ListOrderStatuses lista todos los estados de órdenes de Probability
func (r *Repository) ListOrderStatuses(ctx context.Context, isActive *bool) ([]models.OrderStatus, error) {
	var statuses []models.OrderStatus
	query := r.db.Conn(ctx).Model(&models.OrderStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (r *Repository) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	var mapping models.OrderStatusMapping
	err := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type_id = ? AND original_status = ? AND is_active = ?", integrationTypeID, originalStatus, true).
		Order("priority DESC"). // Si hay múltiples, tomar el de mayor prioridad
		First(&mapping).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró mapeo, retornar nil sin error
			return nil, nil
		}
		return nil, err
	}

	return &mapping.OrderStatusID, nil
}
