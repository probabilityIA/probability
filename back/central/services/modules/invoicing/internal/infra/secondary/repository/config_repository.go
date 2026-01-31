package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type InvoicingConfigRepository struct {
	*Repository
}

func NewInvoicingConfigRepository(repo *Repository) ports.IInvoicingConfigRepository {
	return &InvoicingConfigRepository{Repository: repo}
}

func (r *InvoicingConfigRepository) Create(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	config.ID = model.ID
	return nil
}

func (r *InvoicingConfigRepository) GetByID(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

func (r *InvoicingConfigRepository) GetByIntegration(ctx context.Context, integrationID uint) (*entities.InvoicingConfig, error) {
	var model models.InvoicingConfig

	if err := r.db.Conn(ctx).Where("integration_id = ?", integrationID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("config not found for integration: %w", err)
	}

	return mappers.ConfigToDomain(&model), nil
}

func (r *InvoicingConfigRepository) List(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	var models []*models.InvoicingConfig

	if err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	return mappers.ConfigListToDomain(models), nil
}

func (r *InvoicingConfigRepository) Update(ctx context.Context, config *entities.InvoicingConfig) error {
	model := mappers.ConfigToModel(config)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func (r *InvoicingConfigRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.InvoicingConfig{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}

func (r *InvoicingConfigRepository) ExistsForIntegration(ctx context.Context, integrationID uint) (bool, error) {
	var count int64

	if err := r.db.Conn(ctx).Model(&models.InvoicingConfig{}).
		Where("integration_id = ?", integrationID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check config existence: %w", err)
	}

	return count > 0, nil
}
