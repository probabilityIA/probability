package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
)

func (r *Repository) ListCustomerOrders(ctx context.Context, businessID, userID, customerID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
	base := r.db.Conn(ctx).
		Table("orders o").
		Where("o.business_id = ? AND o.deleted_at IS NULL", businessID)

	if customerID > 0 {
		base = base.Where("(o.user_id = ? OR o.customer_id = ?)", userID, customerID)
	} else {
		base = base.Where("o.user_id = ?", userID)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		ID          string
		OrderNumber string
		TotalAmount float64
		Currency    string
		StatusName  *string
		CreatedAt   string
		ItemsCount  int
	}
	err := base.
		Select(`o.id, o.order_number, o.total_amount, o.currency, os.name AS status_name, o.created_at::text AS created_at,
			(SELECT COUNT(*) FROM order_items oi WHERE oi.order_id = o.id AND oi.deleted_at IS NULL) AS items_count`).
		Joins("LEFT JOIN order_statuses os ON os.id = o.status_id").
		Order("o.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]entities.TiendaOrder, 0, len(rows))
	for _, row := range rows {
		status := ""
		if row.StatusName != nil {
			status = *row.StatusName
		}
		orders = append(orders, entities.TiendaOrder{
			ID:          row.ID,
			OrderNumber: row.OrderNumber,
			Total:       row.TotalAmount,
			Currency:    row.Currency,
			Status:      status,
			CreatedAt:   row.CreatedAt,
			ItemsCount:  row.ItemsCount,
		})
	}
	return orders, total, nil
}
