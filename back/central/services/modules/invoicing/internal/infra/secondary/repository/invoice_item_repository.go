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
type InvoiceItemRepository struct {
	*Repository
}

func NewInvoiceItemRepository(repo *Repository) ports.IInvoiceItemRepository {
	return &InvoiceItemRepository{Repository: repo}
}

func (r *InvoiceItemRepository) Create(ctx context.Context, item *entities.InvoiceItem) error {
=======


func (r *Repository) CreateInvoiceItem(ctx context.Context, item *entities.InvoiceItem) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	model := mappers.InvoiceItemToModel(item)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create invoice item: %w", err)
	}

	item.ID = model.ID
	return nil
}

<<<<<<< HEAD
func (r *InvoiceItemRepository) GetByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error) {
=======
func (r *Repository) GetInvoiceItemsByInvoiceID(ctx context.Context, invoiceID uint) ([]*entities.InvoiceItem, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var models []*models.InvoiceItem

	if err := r.db.Conn(ctx).Where("invoice_id = ?", invoiceID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	return mappers.InvoiceItemListToDomain(models), nil
}

<<<<<<< HEAD
func (r *InvoiceItemRepository) UpdateBatch(ctx context.Context, items []*entities.InvoiceItem) error {
=======
func (r *Repository) UpdateInvoiceItemsBatch(ctx context.Context, items []*entities.InvoiceItem) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	for _, item := range items {
		model := mappers.InvoiceItemToModel(item)
		if err := r.db.Conn(ctx).Save(model).Error; err != nil {
			return fmt.Errorf("failed to update invoice item: %w", err)
		}
	}
	return nil
}
