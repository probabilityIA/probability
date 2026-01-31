package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type InvoiceItemRepository struct {
	*Repository
}

func NewInvoiceItemRepository(repo *Repository) ports.IInvoiceItemRepository {
	return &InvoiceItemRepository{Repository: repo}
}

func (r *InvoiceItemRepository) Create(ctx context.Context, item *entities.InvoiceItem) error {
	model := mappers.InvoiceItemToModel(item)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create invoice item: %w", err)
	}

	item.ID = model.ID
	return nil
}

func (r *InvoiceItemRepository) GetByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error) {
	var models []*models.InvoiceItem

	if err := r.db.Conn(ctx).Where("invoice_id = ?", invoiceID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	return mappers.InvoiceItemListToDomain(models), nil
}

func (r *InvoiceItemRepository) UpdateBatch(ctx context.Context, items []*entities.InvoiceItem) error {
	for _, item := range items {
		model := mappers.InvoiceItemToModel(item)
		if err := r.db.Conn(ctx).Save(model).Error; err != nil {
			return fmt.Errorf("failed to update invoice item: %w", err)
		}
	}
	return nil
}
