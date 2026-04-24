package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

type MockUseCase struct {
	CreateWarehouseFn func(ctx context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error)
	GetWarehouseFn    func(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error)
	ListWarehousesFn  func(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error)
	UpdateWarehouseFn func(ctx context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error)
	DeleteWarehouseFn func(ctx context.Context, businessID, warehouseID uint) error

	CreateLocationFn func(ctx context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error)
	UpdateLocationFn func(ctx context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error)
	DeleteLocationFn func(ctx context.Context, warehouseID, locationID uint, businessID uint) error
	ListLocationsFn  func(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error)

	CreateZoneFn func(ctx context.Context, dto request.CreateZoneDTO) (*entities.WarehouseZone, error)
	GetZoneFn    func(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error)
	ListZonesFn  func(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error)
	UpdateZoneFn func(ctx context.Context, dto request.UpdateZoneDTO) (*entities.WarehouseZone, error)
	DeleteZoneFn func(ctx context.Context, businessID, zoneID uint) error

	CreateAisleFn func(ctx context.Context, dto request.CreateAisleDTO) (*entities.WarehouseAisle, error)
	GetAisleFn    func(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error)
	ListAislesFn  func(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error)
	UpdateAisleFn func(ctx context.Context, dto request.UpdateAisleDTO) (*entities.WarehouseAisle, error)
	DeleteAisleFn func(ctx context.Context, businessID, aisleID uint) error

	CreateRackFn func(ctx context.Context, dto request.CreateRackDTO) (*entities.WarehouseRack, error)
	GetRackFn    func(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error)
	ListRacksFn  func(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error)
	UpdateRackFn func(ctx context.Context, dto request.UpdateRackDTO) (*entities.WarehouseRack, error)
	DeleteRackFn func(ctx context.Context, businessID, rackID uint) error

	CreateRackLevelFn func(ctx context.Context, dto request.CreateRackLevelDTO) (*entities.WarehouseRackLevel, error)
	GetRackLevelFn    func(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error)
	ListRackLevelsFn  func(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error)
	UpdateRackLevelFn func(ctx context.Context, dto request.UpdateRackLevelDTO) (*entities.WarehouseRackLevel, error)
	DeleteRackLevelFn func(ctx context.Context, businessID, levelID uint) error

	GetWarehouseTreeFn func(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error)
}

func (m *MockUseCase) CreateWarehouse(ctx context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error) {
	if m.CreateWarehouseFn != nil {
		return m.CreateWarehouseFn(ctx, dto)
	}
	return &entities.Warehouse{}, nil
}

func (m *MockUseCase) GetWarehouse(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
	if m.GetWarehouseFn != nil {
		return m.GetWarehouseFn(ctx, businessID, warehouseID)
	}
	return &entities.Warehouse{}, nil
}

func (m *MockUseCase) ListWarehouses(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
	if m.ListWarehousesFn != nil {
		return m.ListWarehousesFn(ctx, params)
	}
	return []entities.Warehouse{}, 0, nil
}

func (m *MockUseCase) UpdateWarehouse(ctx context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
	if m.UpdateWarehouseFn != nil {
		return m.UpdateWarehouseFn(ctx, dto)
	}
	return &entities.Warehouse{}, nil
}

func (m *MockUseCase) DeleteWarehouse(ctx context.Context, businessID, warehouseID uint) error {
	if m.DeleteWarehouseFn != nil {
		return m.DeleteWarehouseFn(ctx, businessID, warehouseID)
	}
	return nil
}

func (m *MockUseCase) CreateLocation(ctx context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
	if m.CreateLocationFn != nil {
		return m.CreateLocationFn(ctx, dto)
	}
	return &entities.WarehouseLocation{}, nil
}

func (m *MockUseCase) UpdateLocation(ctx context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
	if m.UpdateLocationFn != nil {
		return m.UpdateLocationFn(ctx, dto)
	}
	return &entities.WarehouseLocation{}, nil
}

func (m *MockUseCase) DeleteLocation(ctx context.Context, warehouseID, locationID uint, businessID uint) error {
	if m.DeleteLocationFn != nil {
		return m.DeleteLocationFn(ctx, warehouseID, locationID, businessID)
	}
	return nil
}

func (m *MockUseCase) ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
	if m.ListLocationsFn != nil {
		return m.ListLocationsFn(ctx, params)
	}
	return []entities.WarehouseLocation{}, nil
}
