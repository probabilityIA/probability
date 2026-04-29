package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) ListWarehouses(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	warehouses, total, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	if len(warehouses) == 0 {
		return warehouses, total, nil
	}

	ids := make([]uint, len(warehouses))
	for i, w := range warehouses {
		ids[i] = w.ID
	}
	counts, err := uc.repo.HierarchyCounts(ctx, ids)
	if err != nil {
		return warehouses, total, nil
	}
	for i := range warehouses {
		c := counts[warehouses[i].ID]
		warehouses[i].ZoneCount = c.Zones
		warehouses[i].AisleCount = c.Aisles
		warehouses[i].RackCount = c.Racks
		warehouses[i].LevelCount = c.Levels
		warehouses[i].PositionCount = c.Positions
	}
	return warehouses, total, nil
}
