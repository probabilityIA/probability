package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/ports"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) HierarchyCounts(ctx context.Context, warehouseIDs []uint) (map[uint]ports.HierarchyCounts, error) {
	out := make(map[uint]ports.HierarchyCounts, len(warehouseIDs))
	if len(warehouseIDs) == 0 {
		return out, nil
	}
	conn := r.db.Conn(ctx)

	type row struct {
		WarehouseID uint
		Total       int
	}

	bump := func(id uint, fn func(*ports.HierarchyCounts)) {
		c := out[id]
		fn(&c)
		out[id] = c
	}

	var zones []row
	if err := conn.Model(&models.WarehouseZone{}).
		Select("warehouse_id, COUNT(*) AS total").
		Where("warehouse_id IN ? AND deleted_at IS NULL", warehouseIDs).
		Group("warehouse_id").
		Scan(&zones).Error; err != nil {
		return nil, err
	}
	for _, z := range zones {
		bump(z.WarehouseID, func(c *ports.HierarchyCounts) { c.Zones = z.Total })
	}

	var aisles []row
	if err := conn.Table("warehouse_aisles AS wa").
		Select("wz.warehouse_id AS warehouse_id, COUNT(*) AS total").
		Joins("JOIN warehouse_zones wz ON wz.id = wa.zone_id AND wz.deleted_at IS NULL").
		Where("wz.warehouse_id IN ? AND wa.deleted_at IS NULL", warehouseIDs).
		Group("wz.warehouse_id").
		Scan(&aisles).Error; err != nil {
		return nil, err
	}
	for _, a := range aisles {
		bump(a.WarehouseID, func(c *ports.HierarchyCounts) { c.Aisles = a.Total })
	}

	var racks []row
	if err := conn.Table("warehouse_racks AS wr").
		Select("wz.warehouse_id AS warehouse_id, COUNT(*) AS total").
		Joins("JOIN warehouse_aisles wa ON wa.id = wr.aisle_id AND wa.deleted_at IS NULL").
		Joins("JOIN warehouse_zones wz ON wz.id = wa.zone_id AND wz.deleted_at IS NULL").
		Where("wz.warehouse_id IN ? AND wr.deleted_at IS NULL", warehouseIDs).
		Group("wz.warehouse_id").
		Scan(&racks).Error; err != nil {
		return nil, err
	}
	for _, rk := range racks {
		bump(rk.WarehouseID, func(c *ports.HierarchyCounts) { c.Racks = rk.Total })
	}

	var levels []row
	if err := conn.Table("warehouse_rack_levels AS wl").
		Select("wz.warehouse_id AS warehouse_id, COUNT(*) AS total").
		Joins("JOIN warehouse_racks wr ON wr.id = wl.rack_id AND wr.deleted_at IS NULL").
		Joins("JOIN warehouse_aisles wa ON wa.id = wr.aisle_id AND wa.deleted_at IS NULL").
		Joins("JOIN warehouse_zones wz ON wz.id = wa.zone_id AND wz.deleted_at IS NULL").
		Where("wz.warehouse_id IN ? AND wl.deleted_at IS NULL", warehouseIDs).
		Group("wz.warehouse_id").
		Scan(&levels).Error; err != nil {
		return nil, err
	}
	for _, lv := range levels {
		bump(lv.WarehouseID, func(c *ports.HierarchyCounts) { c.Levels = lv.Total })
	}

	var positions []row
	if err := conn.Table("warehouse_locations AS loc").
		Select("loc.warehouse_id AS warehouse_id, COUNT(*) AS total").
		Where("loc.warehouse_id IN ? AND loc.deleted_at IS NULL AND loc.level_id IS NOT NULL", warehouseIDs).
		Group("loc.warehouse_id").
		Scan(&positions).Error; err != nil {
		return nil, err
	}
	for _, p := range positions {
		bump(p.WarehouseID, func(c *ports.HierarchyCounts) { c.Positions = p.Total })
	}

	return out, nil
}

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
