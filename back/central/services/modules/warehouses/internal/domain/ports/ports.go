package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	GetByID(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error)
	List(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error)
	Update(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error)
	Delete(ctx context.Context, businessID, warehouseID uint) error
	ExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error)
	ClearDefault(ctx context.Context, businessID uint, excludeID uint) error
	WarehouseExists(ctx context.Context, businessID, warehouseID uint) (bool, error)

	CreateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	GetLocationByID(ctx context.Context, warehouseID, locationID uint) (*entities.WarehouseLocation, error)
	ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error)
	UpdateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error)
	DeleteLocation(ctx context.Context, warehouseID, locationID uint) error
	LocationExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error)

	CreateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error)
	GetZoneByID(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error)
	ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error)
	UpdateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error)
	DeleteZone(ctx context.Context, businessID, zoneID uint) error
	ZoneExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error)

	CreateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error)
	GetAisleByID(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error)
	ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error)
	UpdateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error)
	DeleteAisle(ctx context.Context, businessID, aisleID uint) error
	AisleExistsByCode(ctx context.Context, zoneID uint, code string, excludeID *uint) (bool, error)

	CreateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error)
	GetRackByID(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error)
	ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error)
	UpdateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error)
	DeleteRack(ctx context.Context, businessID, rackID uint) error
	RackExistsByCode(ctx context.Context, aisleID uint, code string, excludeID *uint) (bool, error)

	CreateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error)
	GetRackLevelByID(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error)
	ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error)
	UpdateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error)
	DeleteRackLevel(ctx context.Context, businessID, levelID uint) error
	RackLevelExistsByCode(ctx context.Context, rackID uint, code string, excludeID *uint) (bool, error)

	GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error)
}
