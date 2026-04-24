package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) WarehouseExists(ctx context.Context, businessID, warehouseID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.Warehouse{}).
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", warehouseID, businessID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error) {
	tree := &dtos.WarehouseTreeDTO{WarehouseID: warehouseID}

	var zones []models.WarehouseZone
	if err := r.db.Conn(ctx).
		Where("business_id = ? AND warehouse_id = ? AND deleted_at IS NULL", businessID, warehouseID).
		Order("code ASC").
		Find(&zones).Error; err != nil {
		return nil, err
	}
	if len(zones) == 0 {
		return tree, nil
	}

	zoneIDs := make([]uint, 0, len(zones))
	for i := range zones {
		zoneIDs = append(zoneIDs, zones[i].ID)
	}

	var aisles []models.WarehouseAisle
	if err := r.db.Conn(ctx).
		Where("business_id = ? AND zone_id IN ? AND deleted_at IS NULL", businessID, zoneIDs).
		Order("code ASC").
		Find(&aisles).Error; err != nil {
		return nil, err
	}

	aislesByZone := make(map[uint][]models.WarehouseAisle, len(zoneIDs))
	aisleIDs := make([]uint, 0, len(aisles))
	for i := range aisles {
		aislesByZone[aisles[i].ZoneID] = append(aislesByZone[aisles[i].ZoneID], aisles[i])
		aisleIDs = append(aisleIDs, aisles[i].ID)
	}

	var racks []models.WarehouseRack
	if len(aisleIDs) > 0 {
		if err := r.db.Conn(ctx).
			Where("business_id = ? AND aisle_id IN ? AND deleted_at IS NULL", businessID, aisleIDs).
			Order("code ASC").
			Find(&racks).Error; err != nil {
			return nil, err
		}
	}

	racksByAisle := make(map[uint][]models.WarehouseRack, len(aisleIDs))
	rackIDs := make([]uint, 0, len(racks))
	for i := range racks {
		racksByAisle[racks[i].AisleID] = append(racksByAisle[racks[i].AisleID], racks[i])
		rackIDs = append(rackIDs, racks[i].ID)
	}

	var levels []models.WarehouseRackLevel
	if len(rackIDs) > 0 {
		if err := r.db.Conn(ctx).
			Where("business_id = ? AND rack_id IN ? AND deleted_at IS NULL", businessID, rackIDs).
			Order("ordinal ASC, code ASC").
			Find(&levels).Error; err != nil {
			return nil, err
		}
	}

	levelsByRack := make(map[uint][]models.WarehouseRackLevel, len(rackIDs))
	levelIDs := make([]uint, 0, len(levels))
	for i := range levels {
		levelsByRack[levels[i].RackID] = append(levelsByRack[levels[i].RackID], levels[i])
		levelIDs = append(levelIDs, levels[i].ID)
	}

	var positions []models.WarehouseLocation
	if len(levelIDs) > 0 {
		if err := r.db.Conn(ctx).
			Where("level_id IN ? AND deleted_at IS NULL", levelIDs).
			Order("priority DESC, code ASC").
			Find(&positions).Error; err != nil {
			return nil, err
		}
	}

	positionsByLevel := make(map[uint][]models.WarehouseLocation, len(levelIDs))
	for i := range positions {
		if positions[i].LevelID == nil {
			continue
		}
		positionsByLevel[*positions[i].LevelID] = append(positionsByLevel[*positions[i].LevelID], positions[i])
	}

	tree.Zones = make([]dtos.ZoneNode, 0, len(zones))
	for i := range zones {
		zoneNode := dtos.ZoneNode{WarehouseZone: *zoneModelToEntity(&zones[i])}
		for _, a := range aislesByZone[zones[i].ID] {
			aisleNode := dtos.AisleNode{WarehouseAisle: *aisleModelToEntity(&a)}
			for _, rk := range racksByAisle[a.ID] {
				rackNode := dtos.RackNode{WarehouseRack: *rackModelToEntity(&rk)}
				for _, lv := range levelsByRack[rk.ID] {
					levelNode := dtos.RackLevelNode{WarehouseRackLevel: *rackLevelModelToEntity(&lv)}
					for _, p := range positionsByLevel[lv.ID] {
						levelNode.Positions = append(levelNode.Positions, dtos.PositionNode{WarehouseLocation: *locationModelToEntity(&p)})
					}
					rackNode.Levels = append(rackNode.Levels, levelNode)
				}
				aisleNode.Racks = append(aisleNode.Racks, rackNode)
			}
			zoneNode.Aisles = append(zoneNode.Aisles, aisleNode)
		}
		tree.Zones = append(tree.Zones, zoneNode)
	}

	return tree, nil
}
