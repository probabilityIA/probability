package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

// IRepository define los métodos del repositorio del módulo warehouses
type IRepository interface {
	// Warehouses
	Create(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	GetByID(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error)
	List(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error)
	Update(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	Delete(ctx context.Context, businessID, warehouseID uint) error
	ExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error)
	ClearDefault(ctx context.Context, businessID uint, excludeID uint) error

	// Locations
	CreateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	GetLocationByID(ctx context.Context, warehouseID, locationID uint) (*entities.WarehouseLocation, error)
	ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error)
	UpdateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	DeleteLocation(ctx context.Context, warehouseID, locationID uint) error
	LocationExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error)
}
