package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)



func (r *Repository) CreateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error {
	model := mappers.ProviderToModel(provider)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	provider.ID = model.ID
	return nil
}

func (r *Repository) GetInvoicingProviderByID(ctx context.Context, id uint) (*entities.InvoicingProvider, error) {
	var model models.InvoicingProvider

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

func (r *Repository) GetProviderByBusinessAndType(ctx context.Context, businessID uint, providerTypeCode string) (*entities.InvoicingProvider, error) {
	var model models.InvoicingProvider

	if err := r.db.Conn(ctx).
		Joins("JOIN invoicing_provider_types ON invoicing_provider_types.id = invoicing_providers.provider_type_id").
		Where("invoicing_providers.business_id = ? AND invoicing_provider_types.code = ?", businessID, providerTypeCode).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

func (r *Repository) GetDefaultProviderByBusiness(ctx context.Context, businessID uint) (*entities.InvoicingProvider, error) {
	var model models.InvoicingProvider

	if err := r.db.Conn(ctx).
		Where("business_id = ? AND is_default = ? AND is_active = ?", businessID, true, true).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("default provider not found: %w", err)
	}

	return mappers.ProviderToDomain(&model), nil
}

func (r *Repository) ListInvoicingProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error) {
	var models []*models.InvoicingProvider

	if err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("is_default DESC, created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	return mappers.ProviderListToDomain(models), nil
}

func (r *Repository) UpdateInvoicingProvider(ctx context.Context, provider *entities.InvoicingProvider) error {
	model := mappers.ProviderToModel(provider)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	return nil
}

func (r *Repository) DeleteInvoicingProvider(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.InvoicingProvider{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	return nil
}
