package businessdata

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const maxStatusRows = 15

type totalsRow struct {
	Orders int64
	Amount float64
	Cod    int64
}

type statusRow struct {
	Name  string
	Total int64
}

func (r *Reader) SummarizeOrders(ctx context.Context, businessID uint, from, to time.Time) (*entities.OrdersOverview, error) {
	inPeriod := func() *gorm.DB {
		return r.businessOrders(ctx, businessID).
			Where("orders.created_at >= ? AND orders.created_at < ?", from, to)
	}

	var totals totalsRow
	err := inPeriod().
		Select("COUNT(*) AS orders, COALESCE(SUM(orders.total_amount), 0) AS amount, COUNT(*) FILTER (WHERE orders.is_cod) AS cod").
		Scan(&totals).Error
	if err != nil {
		return nil, err
	}

	overview := &entities.OrdersOverview{
		From:           from,
		To:             to,
		Orders:         totals.Orders,
		Amount:         totals.Amount,
		CashOnDelivery: totals.Cod,
	}

	if err := inPeriod().Where(withoutGuideSQL).Count(&overview.WithoutGuide).Error; err != nil {
		return nil, err
	}
	if err := inPeriod().Where(withNoveltySQL).Count(&overview.WithNovelty).Error; err != nil {
		return nil, err
	}

	var byStatus []statusRow
	err = inPeriod().
		Joins("LEFT JOIN order_statuses ON order_statuses.id = orders.status_id").
		Select("COALESCE(order_statuses.name, orders.status) AS name, COUNT(*) AS total").
		Group("COALESCE(order_statuses.name, orders.status)").
		Order("total DESC").
		Limit(maxStatusRows).
		Scan(&byStatus).Error
	if err != nil {
		return nil, err
	}
	for _, row := range byStatus {
		overview.ByStatus = append(overview.ByStatus, entities.StatusCount{Name: row.Name, Count: row.Total})
	}

	var shipments []statusRow
	err = r.db.Conn(ctx).Model(&models.Shipment{}).
		Joins("JOIN orders ON orders.id = shipments.order_id").
		Where("orders.business_id = ? AND orders.deleted_at IS NULL", businessID).
		Where("shipments.created_at >= ? AND shipments.created_at < ?", from, to).
		Select("shipments.status AS name, COUNT(*) AS total").
		Group("shipments.status").
		Order("total DESC").
		Scan(&shipments).Error
	if err != nil {
		return nil, err
	}
	for _, row := range shipments {
		overview.ShipmentStatus = append(overview.ShipmentStatus, entities.StatusCount{Name: row.Name, Count: row.Total})
	}
	return overview, nil
}
