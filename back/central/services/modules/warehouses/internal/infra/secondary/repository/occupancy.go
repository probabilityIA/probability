package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (r *Repository) GetWarehouseOccupancy(ctx context.Context, businessID, warehouseID uint) ([]entities.OccupancyItem, error) {
	var rows []struct {
		LocationID uint
		Quantity   int
		Reserved   int
		Capacity   *int
	}

	err := r.db.Conn(ctx).
		Table("warehouse_locations AS l").
		Select("l.id AS location_id, COALESCE(SUM(il.quantity),0) AS quantity, COALESCE(SUM(il.reserved_qty),0) AS reserved, l.capacity AS capacity").
		Joins("LEFT JOIN inventory_levels il ON il.location_id = l.id AND il.deleted_at IS NULL AND il.business_id = ?", businessID).
		Where("l.warehouse_id = ? AND l.deleted_at IS NULL", warehouseID).
		Group("l.id, l.capacity").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]entities.OccupancyItem, len(rows))
	for i, row := range rows {
		items[i] = entities.OccupancyItem{
			LocationID: row.LocationID,
			Quantity:   row.Quantity,
			Reserved:   row.Reserved,
			Capacity:   row.Capacity,
		}
	}
	return items, nil
}
