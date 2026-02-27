// Package mocks contiene los mocks manuales para el módulo warehouses.
// Los mocks se organizan aquí para ser reutilizables en todos los archivos de test.
package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

// RepositoryMock es un mock manual de ports.IRepository.
// Cada método expone un campo Fn para que el test inyecte el comportamiento deseado.
// Si el campo Fn es nil, el método retorna el valor cero correspondiente sin error.
type RepositoryMock struct {
	// --- Warehouses ---
	CreateFn         func(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	GetByIDFn        func(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error)
	ListFn           func(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error)
	UpdateFn         func(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	DeleteFn         func(ctx context.Context, businessID, warehouseID uint) error
	ExistsByCodeFn   func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error)
	ClearDefaultFn   func(ctx context.Context, businessID uint, excludeID uint) error

	// --- Locations ---
	CreateLocationFn        func(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	GetLocationByIDFn       func(ctx context.Context, warehouseID, locationID uint) (*entities.WarehouseLocation, error)
	ListLocationsFn         func(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error)
	UpdateLocationFn        func(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	DeleteLocationFn        func(ctx context.Context, warehouseID, locationID uint) error
	LocationExistsByCodeFn  func(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error)
}

// --- Implementación de ports.IRepository ---

func (m *RepositoryMock) Create(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error) {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, warehouse)
	}
	return warehouse, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, businessID, warehouseID)
	}
	return &entities.Warehouse{ID: warehouseID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return []entities.Warehouse{}, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, warehouse)
	}
	return warehouse, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, businessID, warehouseID uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, businessID, warehouseID)
	}
	return nil
}

func (m *RepositoryMock) ExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
	if m.ExistsByCodeFn != nil {
		return m.ExistsByCodeFn(ctx, businessID, code, excludeID)
	}
	return false, nil
}

func (m *RepositoryMock) ClearDefault(ctx context.Context, businessID uint, excludeID uint) error {
	if m.ClearDefaultFn != nil {
		return m.ClearDefaultFn(ctx, businessID, excludeID)
	}
	return nil
}

func (m *RepositoryMock) CreateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
	if m.CreateLocationFn != nil {
		return m.CreateLocationFn(ctx, location)
	}
	return location, nil
}

func (m *RepositoryMock) GetLocationByID(ctx context.Context, warehouseID, locationID uint) (*entities.WarehouseLocation, error) {
	if m.GetLocationByIDFn != nil {
		return m.GetLocationByIDFn(ctx, warehouseID, locationID)
	}
	return &entities.WarehouseLocation{ID: locationID, WarehouseID: warehouseID}, nil
}

func (m *RepositoryMock) ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
	if m.ListLocationsFn != nil {
		return m.ListLocationsFn(ctx, params)
	}
	return []entities.WarehouseLocation{}, nil
}

func (m *RepositoryMock) UpdateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
	if m.UpdateLocationFn != nil {
		return m.UpdateLocationFn(ctx, location)
	}
	return location, nil
}

func (m *RepositoryMock) DeleteLocation(ctx context.Context, warehouseID, locationID uint) error {
	if m.DeleteLocationFn != nil {
		return m.DeleteLocationFn(ctx, warehouseID, locationID)
	}
	return nil
}

func (m *RepositoryMock) LocationExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error) {
	if m.LocationExistsByCodeFn != nil {
		return m.LocationExistsByCodeFn(ctx, warehouseID, code, excludeID)
	}
	return false, nil
}
