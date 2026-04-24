package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (m *MockUseCase) CreateZone(ctx context.Context, dto request.CreateZoneDTO) (*entities.WarehouseZone, error) {
	if m.CreateZoneFn != nil {
		return m.CreateZoneFn(ctx, dto)
	}
	return &entities.WarehouseZone{}, nil
}

func (m *MockUseCase) GetZone(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error) {
	if m.GetZoneFn != nil {
		return m.GetZoneFn(ctx, businessID, zoneID)
	}
	return &entities.WarehouseZone{ID: zoneID, BusinessID: businessID}, nil
}

func (m *MockUseCase) ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error) {
	if m.ListZonesFn != nil {
		return m.ListZonesFn(ctx, params)
	}
	return []entities.WarehouseZone{}, 0, nil
}

func (m *MockUseCase) UpdateZone(ctx context.Context, dto request.UpdateZoneDTO) (*entities.WarehouseZone, error) {
	if m.UpdateZoneFn != nil {
		return m.UpdateZoneFn(ctx, dto)
	}
	return &entities.WarehouseZone{}, nil
}

func (m *MockUseCase) DeleteZone(ctx context.Context, businessID, zoneID uint) error {
	if m.DeleteZoneFn != nil {
		return m.DeleteZoneFn(ctx, businessID, zoneID)
	}
	return nil
}

func (m *MockUseCase) CreateAisle(ctx context.Context, dto request.CreateAisleDTO) (*entities.WarehouseAisle, error) {
	if m.CreateAisleFn != nil {
		return m.CreateAisleFn(ctx, dto)
	}
	return &entities.WarehouseAisle{}, nil
}

func (m *MockUseCase) GetAisle(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error) {
	if m.GetAisleFn != nil {
		return m.GetAisleFn(ctx, businessID, aisleID)
	}
	return &entities.WarehouseAisle{ID: aisleID, BusinessID: businessID}, nil
}

func (m *MockUseCase) ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error) {
	if m.ListAislesFn != nil {
		return m.ListAislesFn(ctx, params)
	}
	return []entities.WarehouseAisle{}, 0, nil
}

func (m *MockUseCase) UpdateAisle(ctx context.Context, dto request.UpdateAisleDTO) (*entities.WarehouseAisle, error) {
	if m.UpdateAisleFn != nil {
		return m.UpdateAisleFn(ctx, dto)
	}
	return &entities.WarehouseAisle{}, nil
}

func (m *MockUseCase) DeleteAisle(ctx context.Context, businessID, aisleID uint) error {
	if m.DeleteAisleFn != nil {
		return m.DeleteAisleFn(ctx, businessID, aisleID)
	}
	return nil
}

func (m *MockUseCase) CreateRack(ctx context.Context, dto request.CreateRackDTO) (*entities.WarehouseRack, error) {
	if m.CreateRackFn != nil {
		return m.CreateRackFn(ctx, dto)
	}
	return &entities.WarehouseRack{}, nil
}

func (m *MockUseCase) GetRack(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error) {
	if m.GetRackFn != nil {
		return m.GetRackFn(ctx, businessID, rackID)
	}
	return &entities.WarehouseRack{ID: rackID, BusinessID: businessID}, nil
}

func (m *MockUseCase) ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error) {
	if m.ListRacksFn != nil {
		return m.ListRacksFn(ctx, params)
	}
	return []entities.WarehouseRack{}, 0, nil
}

func (m *MockUseCase) UpdateRack(ctx context.Context, dto request.UpdateRackDTO) (*entities.WarehouseRack, error) {
	if m.UpdateRackFn != nil {
		return m.UpdateRackFn(ctx, dto)
	}
	return &entities.WarehouseRack{}, nil
}

func (m *MockUseCase) DeleteRack(ctx context.Context, businessID, rackID uint) error {
	if m.DeleteRackFn != nil {
		return m.DeleteRackFn(ctx, businessID, rackID)
	}
	return nil
}

func (m *MockUseCase) CreateRackLevel(ctx context.Context, dto request.CreateRackLevelDTO) (*entities.WarehouseRackLevel, error) {
	if m.CreateRackLevelFn != nil {
		return m.CreateRackLevelFn(ctx, dto)
	}
	return &entities.WarehouseRackLevel{}, nil
}

func (m *MockUseCase) GetRackLevel(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error) {
	if m.GetRackLevelFn != nil {
		return m.GetRackLevelFn(ctx, businessID, levelID)
	}
	return &entities.WarehouseRackLevel{ID: levelID, BusinessID: businessID}, nil
}

func (m *MockUseCase) ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error) {
	if m.ListRackLevelsFn != nil {
		return m.ListRackLevelsFn(ctx, params)
	}
	return []entities.WarehouseRackLevel{}, 0, nil
}

func (m *MockUseCase) UpdateRackLevel(ctx context.Context, dto request.UpdateRackLevelDTO) (*entities.WarehouseRackLevel, error) {
	if m.UpdateRackLevelFn != nil {
		return m.UpdateRackLevelFn(ctx, dto)
	}
	return &entities.WarehouseRackLevel{}, nil
}

func (m *MockUseCase) DeleteRackLevel(ctx context.Context, businessID, levelID uint) error {
	if m.DeleteRackLevelFn != nil {
		return m.DeleteRackLevelFn(ctx, businessID, levelID)
	}
	return nil
}

func (m *MockUseCase) GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error) {
	if m.GetWarehouseTreeFn != nil {
		return m.GetWarehouseTreeFn(ctx, businessID, warehouseID)
	}
	return &dtos.WarehouseTreeDTO{WarehouseID: warehouseID}, nil
}
