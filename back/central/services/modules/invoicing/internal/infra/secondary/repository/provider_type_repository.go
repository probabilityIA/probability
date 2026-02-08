package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)



func (r *Repository) GetProviderTypeByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error) {
	var model models.InvoicingProviderType

	if err := r.db.Conn(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		return nil, fmt.Errorf("provider type not found: %w", err)
	}

	return mappers.ProviderTypeToDomain(&model), nil
}

func (r *Repository) ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
	var models []*models.InvoicingProviderType

	if err := r.db.Conn(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}

func (r *Repository) GetActiveProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
	var models []*models.InvoicingProviderType

	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list active provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}
