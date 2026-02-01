package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

type providerTypeRepository struct {
	*Repository
}

// NewProviderTypeRepository crea una nueva instancia del repositorio de tipos de proveedor
func NewProviderTypeRepository(repo *Repository) ports.IProviderTypeRepository {
	return &providerTypeRepository{Repository: repo}
}

// GetByCode obtiene un tipo de proveedor por código
func (r *providerTypeRepository) GetByCode(ctx context.Context, code string) (*entities.ProviderType, error) {
	var model models.ProviderType

	if err := r.db.Conn(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("provider type not found")
		}
		r.log.Error(ctx).Err(err).Str("code", code).Msg("Failed to get provider type")
		return nil, fmt.Errorf("failed to get provider type: %w", err)
	}

	return mappers.ProviderTypeToDomain(&model), nil
}

// List lista todos los tipos de proveedores
func (r *providerTypeRepository) List(ctx context.Context) ([]*entities.ProviderType, error) {
	var models []*models.ProviderType

	if err := r.db.Conn(ctx).Order("name ASC").Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to list provider types")
		return nil, fmt.Errorf("failed to list provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}

// GetActive lista los tipos de proveedores activos
func (r *providerTypeRepository) GetActive(ctx context.Context) ([]*entities.ProviderType, error) {
	var models []*models.ProviderType

	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to list active provider types")
		return nil, fmt.Errorf("failed to list active provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}
