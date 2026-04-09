package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) UpsertCustomerOrderItem(ctx context.Context, item *entities.CustomerOrderItem) error {
	var existing models.CustomerOrderItem
	err := r.db.Conn(ctx).
		Where("customer_id = ? AND order_id = ? AND product_id = ?",
			item.CustomerID, item.OrderID, item.ProductID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		model := mapCustomerOrderItemFromEntity(item)
		return r.db.Conn(ctx).Create(model).Error
	}
	return err
}

func (r *Repository) UpdateOrderItemsStatus(ctx context.Context, orderID string, status string) error {
	return r.db.Conn(ctx).Model(&models.CustomerOrderItem{}).
		Where("order_id = ?", orderID).
		Update("order_status", status).Error
}
