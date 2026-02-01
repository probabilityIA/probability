package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

type providerRepository struct {
	*Repository
}

// NewProviderRepository crea una nueva instancia del repositorio de proveedores
func NewProviderRepository(repo *Repository) ports.IProviderRepository {
	return &providerRepository{Repository: repo}
}

// Create crea un nuevo proveedor en la base de datos
func (r *providerRepository) Create(ctx context.Context, provider *entities.Provider) error {
	model := mappers.ProviderToModel(provider)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to create provider")
		return fmt.Errorf("failed to create provider: %w", err)
	}

	provider.ID = model.ID
	provider.CreatedAt = model.CreatedAt
	provider.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID obtiene un proveedor por ID
func (r *providerRepository) GetByID(ctx context.Context, id uint) (*entities.Provider, error) {
	var model models.Provider

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("provider not found")
		}
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Failed to get provider")
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

// GetByBusinessAndType obtiene un proveedor por negocio y tipo de proveedor
func (r *providerRepository) GetByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.Provider, error) {
	var model models.Provider

	err := r.db.Conn(ctx).
		Joins("JOIN invoicing_provider_types ON invoicing_provider_types.id = invoicing_providers.provider_type_id").
		Where("invoicing_providers.business_id = ? AND invoicing_provider_types.code = ?", businessID, providerTypeCode).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No error, solo no encontrado
		}
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Str("type_code", providerTypeCode).Msg("Failed to get provider")
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

// GetDefaultByBusiness obtiene el proveedor por defecto de un negocio
func (r *providerRepository) GetDefaultByBusiness(ctx context.Context, businessID uint) (*entities.Provider, error) {
	var model models.Provider

	err := r.db.Conn(ctx).
		Where("business_id = ? AND is_default = ? AND is_active = ?", businessID, true, true).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No error, solo no encontrado
		}
		r.log.Error(ctx).Err(err).Uint("business_id", businessID).Msg("Failed to get default provider")
		return nil, fmt.Errorf("failed to get default provider: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

// List lista proveedores según filtros
func (r *providerRepository) List(ctx context.Context, filters *dtos.ProviderFiltersDTO) ([]*entities.Provider, error) {
	var models []*models.Provider

	query := r.db.Conn(ctx)

	// Aplicar filtros
	if filters.BusinessID != nil {
		query = query.Where("business_id = ?", *filters.BusinessID)
	}

	if filters.ProviderTypeCode != nil {
		query = query.Joins("JOIN invoicing_provider_types ON invoicing_provider_types.id = invoicing_providers.provider_type_id").
			Where("invoicing_provider_types.code = ?", *filters.ProviderTypeCode)
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}

	// Orden por defecto
	query = query.Order("is_default DESC, created_at DESC")

	// Paginación
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to list providers")
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	return mappers.ProviderListToDomain(models), nil
}

// Update actualiza un proveedor existente
func (r *providerRepository) Update(ctx context.Context, provider *entities.Provider) error {
	model := mappers.ProviderToModel(provider)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", provider.ID).Msg("Failed to update provider")
		return fmt.Errorf("failed to update provider: %w", err)
	}

	provider.UpdatedAt = model.UpdatedAt

	return nil
}

// Delete elimina un proveedor (soft delete)
func (r *providerRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.Provider{}, id).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("id", id).Msg("Failed to delete provider")
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	return nil
}
