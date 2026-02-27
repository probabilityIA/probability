package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

// MockUseCase es un mock inyectable de app.IUseCase para tests unitarios.
// Cada campo con sufijo Fn permite configurar el comportamiento por test.
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
