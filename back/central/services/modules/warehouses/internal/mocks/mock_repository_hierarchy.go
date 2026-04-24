package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (m *RepositoryMock) WarehouseExists(ctx context.Context, businessID, warehouseID uint) (bool, error) {
	if m.WarehouseExistsFn != nil {
		return m.WarehouseExistsFn(ctx, businessID, warehouseID)
	}
	return true, nil
}

func (m *RepositoryMock) CreateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error) {
	if m.CreateZoneFn != nil {
		return m.CreateZoneFn(ctx, zone)
	}
	return zone, nil
}

func (m *RepositoryMock) GetZoneByID(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error) {
	if m.GetZoneByIDFn != nil {
		return m.GetZoneByIDFn(ctx, businessID, zoneID)
	}
	return &entities.WarehouseZone{ID: zoneID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error) {
	if m.ListZonesFn != nil {
		return m.ListZonesFn(ctx, params)
	}
	return []entities.WarehouseZone{}, 0, nil
}

func (m *RepositoryMock) UpdateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error) {
	if m.UpdateZoneFn != nil {
		return m.UpdateZoneFn(ctx, zone)
	}
	return zone, nil
}

func (m *RepositoryMock) DeleteZone(ctx context.Context, businessID, zoneID uint) error {
	if m.DeleteZoneFn != nil {
		return m.DeleteZoneFn(ctx, businessID, zoneID)
	}
	return nil
}

func (m *RepositoryMock) ZoneExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error) {
	if m.ZoneExistsByCodeFn != nil {
		return m.ZoneExistsByCodeFn(ctx, warehouseID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error) {
	if m.CreateAisleFn != nil {
		return m.CreateAisleFn(ctx, aisle)
	}
	return aisle, nil
}

func (m *RepositoryMock) GetAisleByID(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error) {
	if m.GetAisleByIDFn != nil {
		return m.GetAisleByIDFn(ctx, businessID, aisleID)
	}
	return &entities.WarehouseAisle{ID: aisleID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error) {
	if m.ListAislesFn != nil {
		return m.ListAislesFn(ctx, params)
	}
	return []entities.WarehouseAisle{}, 0, nil
}

func (m *RepositoryMock) UpdateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error) {
	if m.UpdateAisleFn != nil {
		return m.UpdateAisleFn(ctx, aisle)
	}
	return aisle, nil
}

func (m *RepositoryMock) DeleteAisle(ctx context.Context, businessID, aisleID uint) error {
	if m.DeleteAisleFn != nil {
		return m.DeleteAisleFn(ctx, businessID, aisleID)
	}
	return nil
}

func (m *RepositoryMock) AisleExistsByCode(ctx context.Context, zoneID uint, code string, excludeID *uint) (bool, error) {
	if m.AisleExistsByCodeFn != nil {
		return m.AisleExistsByCodeFn(ctx, zoneID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error) {
	if m.CreateRackFn != nil {
		return m.CreateRackFn(ctx, rack)
	}
	return rack, nil
}

func (m *RepositoryMock) GetRackByID(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error) {
	if m.GetRackByIDFn != nil {
		return m.GetRackByIDFn(ctx, businessID, rackID)
	}
	return &entities.WarehouseRack{ID: rackID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error) {
	if m.ListRacksFn != nil {
		return m.ListRacksFn(ctx, params)
	}
	return []entities.WarehouseRack{}, 0, nil
}

func (m *RepositoryMock) UpdateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error) {
	if m.UpdateRackFn != nil {
		return m.UpdateRackFn(ctx, rack)
	}
	return rack, nil
}

func (m *RepositoryMock) DeleteRack(ctx context.Context, businessID, rackID uint) error {
	if m.DeleteRackFn != nil {
		return m.DeleteRackFn(ctx, businessID, rackID)
	}
	return nil
}

func (m *RepositoryMock) RackExistsByCode(ctx context.Context, aisleID uint, code string, excludeID *uint) (bool, error) {
	if m.RackExistsByCodeFn != nil {
		return m.RackExistsByCodeFn(ctx, aisleID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error) {
	if m.CreateRackLevelFn != nil {
		return m.CreateRackLevelFn(ctx, level)
	}
	return level, nil
}

func (m *RepositoryMock) GetRackLevelByID(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error) {
	if m.GetRackLevelByIDFn != nil {
		return m.GetRackLevelByIDFn(ctx, businessID, levelID)
	}
	return &entities.WarehouseRackLevel{ID: levelID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error) {
	if m.ListRackLevelsFn != nil {
		return m.ListRackLevelsFn(ctx, params)
	}
	return []entities.WarehouseRackLevel{}, 0, nil
}

func (m *RepositoryMock) UpdateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error) {
	if m.UpdateRackLevelFn != nil {
		return m.UpdateRackLevelFn(ctx, level)
	}
	return level, nil
}

func (m *RepositoryMock) DeleteRackLevel(ctx context.Context, businessID, levelID uint) error {
	if m.DeleteRackLevelFn != nil {
		return m.DeleteRackLevelFn(ctx, businessID, levelID)
	}
	return nil
}

func (m *RepositoryMock) RackLevelExistsByCode(ctx context.Context, rackID uint, code string, excludeID *uint) (bool, error) {
	if m.RackLevelExistsByCodeFn != nil {
		return m.RackLevelExistsByCodeFn(ctx, rackID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error) {
	if m.GetWarehouseTreeFn != nil {
		return m.GetWarehouseTreeFn(ctx, businessID, warehouseID)
	}
	return &dtos.WarehouseTreeDTO{WarehouseID: warehouseID}, nil
}
