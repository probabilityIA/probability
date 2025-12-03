package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/payments/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ═══════════════════════════════════════════
// CONSTRUCTOR
// ═══════════════════════════════════════════

// Repository implementa domain.IRepository
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

// ═══════════════════════════════════════════
// PAYMENT METHOD REPOSITORY
// ═══════════════════════════════════════════

func (r *Repository) CreatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error {
	return r.db.Conn(ctx).Create(method).Error
}

func (r *Repository) GetPaymentMethodByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.Conn(ctx).Where("id = ?", id).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment method not found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *Repository) GetPaymentMethodByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	err := r.db.Conn(ctx).Where("code = ?", code).First(&method).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment method not found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *Repository) ListPaymentMethods(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]models.PaymentMethod, int64, error) {
	var methods []models.PaymentMethod
	var total int64

	query := r.db.Conn(ctx).Model(&models.PaymentMethod{})

	// Aplicar filtros
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginar
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&methods).Error; err != nil {
		return nil, 0, err
	}

	return methods, total, nil
}

func (r *Repository) UpdatePaymentMethod(ctx context.Context, method *models.PaymentMethod) error {
	return r.db.Conn(ctx).Save(method).Error
}

func (r *Repository) DeletePaymentMethod(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.PaymentMethod{}, id).Error
}

func (r *Repository) TogglePaymentMethodActive(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	method, err := r.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	method.IsActive = !method.IsActive
	if err := r.UpdatePaymentMethod(ctx, method); err != nil {
		return nil, err
	}

	return method, nil
}

func (r *Repository) PaymentMethodExists(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.PaymentMethod{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *Repository) PaymentMethodHasActiveMappings(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.PaymentMethodMapping{}).
		Where("payment_method_id = ? AND is_active = ?", id, true).
		Count(&count).Error
	return count > 0, err
}

// ═══════════════════════════════════════════
// PAYMENT MAPPING REPOSITORY
// ═══════════════════════════════════════════

func (r *Repository) CreatePaymentMapping(ctx context.Context, mapping *models.PaymentMethodMapping) error {
	return r.db.Conn(ctx).Create(mapping).Error
}

func (r *Repository) GetPaymentMappingByID(ctx context.Context, id uint) (*models.PaymentMethodMapping, error) {
	var mapping models.PaymentMethodMapping
	err := r.db.Conn(ctx).Where("id = ?", id).First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment mapping not found")
		}
		return nil, err
	}
	return &mapping, nil
}

func (r *Repository) GetPaymentMappingByIDWithMethod(ctx context.Context, id uint) (*models.PaymentMethodMapping, error) {
	var mapping models.PaymentMethodMapping
	err := r.db.Conn(ctx).Preload("PaymentMethod").Where("id = ?", id).First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment mapping not found")
		}
		return nil, err
	}
	return &mapping, nil
}

func (r *Repository) ListPaymentMappings(ctx context.Context, filters map[string]interface{}) ([]models.PaymentMethodMapping, int64, error) {
	var mappings []models.PaymentMethodMapping
	var total int64

	query := r.db.Conn(ctx).Model(&models.PaymentMethodMapping{})

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

func (r *Repository) ListPaymentMappingsWithMethods(ctx context.Context, filters map[string]interface{}) ([]models.PaymentMethodMapping, int64, error) {
	var mappings []models.PaymentMethodMapping
	var total int64

	query := r.db.Conn(ctx).Model(&models.PaymentMethodMapping{}).Preload("PaymentMethod")

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

func (r *Repository) UpdatePaymentMapping(ctx context.Context, mapping *models.PaymentMethodMapping) error {
	return r.db.Conn(ctx).Save(mapping).Error
}

func (r *Repository) DeletePaymentMapping(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.PaymentMethodMapping{}, id).Error
}

func (r *Repository) GetPaymentMappingsByIntegrationType(ctx context.Context, integrationType string) ([]models.PaymentMethodMapping, error) {
	var mappings []models.PaymentMethodMapping
	err := r.db.Conn(ctx).
		Where("integration_type = ?", integrationType).
		Order("priority DESC, created_at DESC").
		Find(&mappings).Error
	return mappings, err
}

func (r *Repository) GetPaymentMappingsByIntegrationTypeWithMethods(ctx context.Context, integrationType string) ([]models.PaymentMethodMapping, error) {
	var mappings []models.PaymentMethodMapping
	err := r.db.Conn(ctx).
		Preload("PaymentMethod").
		Where("integration_type = ?", integrationType).
		Order("priority DESC, created_at DESC").
		Find(&mappings).Error
	return mappings, err
}

func (r *Repository) TogglePaymentMappingActive(ctx context.Context, id uint) (*models.PaymentMethodMapping, error) {
	mapping, err := r.GetPaymentMappingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mapping.IsActive = !mapping.IsActive
	if err := r.UpdatePaymentMapping(ctx, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (r *Repository) PaymentMappingExists(ctx context.Context, integrationType, originalMethod string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.PaymentMethodMapping{}).
		Where("integration_type = ? AND original_method = ?", integrationType, originalMethod).
		Count(&count).Error
	return count > 0, err
}
