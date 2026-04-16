package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) UpsertCustomerProductHistory(ctx context.Context, product *entities.CustomerProductHistory) error {
	var existing models.CustomerProductHistory
	err := r.db.Conn(ctx).
		Where("customer_id = ? AND business_id = ? AND product_id = ?",
			product.CustomerID, product.BusinessID, product.ProductID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		model := mapCustomerProductFromEntity(product)
		return r.db.Conn(ctx).Create(model).Error
	}
	if err != nil {
		return err
	}

	return r.db.Conn(ctx).Model(&existing).Updates(map[string]any{
		"times_ordered":   existing.TimesOrdered + product.TimesOrdered,
		"total_quantity":  existing.TotalQuantity + product.TotalQuantity,
		"total_spent":     existing.TotalSpent + product.TotalSpent,
		"product_name":    coalesceString(product.ProductName, existing.ProductName),
		"product_sku":     coalesceString(product.ProductSKU, existing.ProductSKU),
		"product_image":   coalesceProductImage(product.ProductImage, existing.ProductImage),
		"last_ordered_at": product.LastOrderedAt,
	}).Error
}

func coalesceProductImage(preferred, fallback *string) *string {
	if preferred != nil && *preferred != "" {
		return preferred
	}
	return fallback
}
