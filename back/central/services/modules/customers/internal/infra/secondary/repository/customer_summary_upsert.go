package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) UpsertCustomerSummary(ctx context.Context, summary *entities.CustomerSummary) error {
	var existing models.CustomerSummary
	err := r.db.Conn(ctx).
		Where("customer_id = ? AND business_id = ?", summary.CustomerID, summary.BusinessID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		model := mapCustomerSummaryFromEntity(summary)
		return r.db.Conn(ctx).Create(model).Error
	}
	if err != nil {
		return err
	}

	newTotalOrders := existing.TotalOrders + summary.TotalOrders
	newTotalSpent := existing.TotalSpent + summary.TotalSpent

	return r.db.Conn(ctx).Model(&existing).Updates(map[string]any{
		"total_orders":       newTotalOrders,
		"delivered_orders":   existing.DeliveredOrders + summary.DeliveredOrders,
		"cancelled_orders":   existing.CancelledOrders + summary.CancelledOrders,
		"in_progress_orders": existing.InProgressOrders + summary.InProgressOrders,
		"total_spent":        newTotalSpent,
		"total_paid_orders":  existing.TotalPaidOrders + summary.TotalPaidOrders,
		"avg_ticket":         calculateAvgTicket(newTotalSpent, newTotalOrders),
		"avg_delivery_score": coalesceFloat(summary.AvgDeliveryScore, existing.AvgDeliveryScore),
		"first_order_at":     coalesceTime(existing.FirstOrderAt, summary.FirstOrderAt),
		"last_order_at":      latestTime(existing.LastOrderAt, summary.LastOrderAt),
		"preferred_platform": coalesceString(summary.PreferredPlatform, existing.PreferredPlatform),
		"last_updated_at":    time.Now(),
	}).Error
}
