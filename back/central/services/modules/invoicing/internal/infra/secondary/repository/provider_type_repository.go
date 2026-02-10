package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

<<<<<<< HEAD
type InvoicingProviderTypeRepository struct {
	*Repository
}

func NewInvoicingProviderTypeRepository(repo *Repository) ports.IInvoicingProviderTypeRepository {
	return &InvoicingProviderTypeRepository{Repository: repo}
}

func (r *InvoicingProviderTypeRepository) GetByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error) {
=======


func (r *Repository) GetProviderTypeByCode(ctx context.Context, code string) (*entities.InvoicingProviderType, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.InvoicingProviderType

	if err := r.db.Conn(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		return nil, fmt.Errorf("provider type not found: %w", err)
	}

	return mappers.ProviderTypeToDomain(&model), nil
}

<<<<<<< HEAD
func (r *InvoicingProviderTypeRepository) List(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
=======
func (r *Repository) ListProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var models []*models.InvoicingProviderType

	if err := r.db.Conn(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}

<<<<<<< HEAD
func (r *InvoicingProviderTypeRepository) GetActive(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
=======
func (r *Repository) GetActiveProviderTypes(ctx context.Context) ([]*entities.InvoicingProviderType, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var models []*models.InvoicingProviderType

	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list active provider types: %w", err)
	}

	return mappers.ProviderTypeListToDomain(models), nil
}
