package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/ports"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) HierarchyDepth(ctx context.Context, warehouseID uint) (ports.HierarchyDepth, error) {
	depth := ports.HierarchyDepth{}

	var zoneCount int64
	if err := r.db.Conn(ctx).Model(&models.WarehouseZone{}).
		Where("warehouse_id = ? AND deleted_at IS NULL", warehouseID).
		Count(&zoneCount).Error; err != nil {
		return depth, err
	}
	depth.HasZones = zoneCount > 0

	var rackCount int64
	if err := r.db.Conn(ctx).Model(&models.WarehouseRack{}).
		Joins("JOIN warehouse_aisles wa ON wa.id = warehouse_racks.aisle_id AND wa.deleted_at IS NULL").
		Joins("JOIN warehouse_zones wz ON wz.id = wa.zone_id AND wz.deleted_at IS NULL").
		Where("wz.warehouse_id = ? AND warehouse_racks.deleted_at IS NULL", warehouseID).
		Count(&rackCount).Error; err != nil {
		return depth, err
	}
	depth.HasRacks = rackCount > 0

	return depth, nil
}
