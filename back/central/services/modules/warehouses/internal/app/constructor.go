package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/ports"
)

// IUseCase define los casos de uso del módulo warehouses
type IUseCase interface {
	CreateWarehouse(ctx context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error)
	GetWarehouse(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error)
	ListWarehouses(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error)
	UpdateWarehouse(ctx context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error)
	DeleteWarehouse(ctx context.Context, businessID, warehouseID uint) error

	CreateLocation(ctx context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error)
	UpdateLocation(ctx context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error)
	DeleteLocation(ctx context.Context, warehouseID, locationID uint, businessID uint) error
	ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error)
}

// UseCase implementa IUseCase
type UseCase struct {
	repo ports.IRepository
}

// New crea una nueva instancia del use case
func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
