package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/ports"
)

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

	CreateZone(ctx context.Context, dto request.CreateZoneDTO) (*entities.WarehouseZone, error)
	GetZone(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error)
	ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error)
	UpdateZone(ctx context.Context, dto request.UpdateZoneDTO) (*entities.WarehouseZone, error)
	DeleteZone(ctx context.Context, businessID, zoneID uint) error

	CreateAisle(ctx context.Context, dto request.CreateAisleDTO) (*entities.WarehouseAisle, error)
	GetAisle(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error)
	ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error)
	UpdateAisle(ctx context.Context, dto request.UpdateAisleDTO) (*entities.WarehouseAisle, error)
	DeleteAisle(ctx context.Context, businessID, aisleID uint) error

	CreateRack(ctx context.Context, dto request.CreateRackDTO) (*entities.WarehouseRack, error)
	GetRack(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error)
	ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error)
	UpdateRack(ctx context.Context, dto request.UpdateRackDTO) (*entities.WarehouseRack, error)
	DeleteRack(ctx context.Context, businessID, rackID uint) error

	CreateRackLevel(ctx context.Context, dto request.CreateRackLevelDTO) (*entities.WarehouseRackLevel, error)
	GetRackLevel(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error)
	ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error)
	UpdateRackLevel(ctx context.Context, dto request.UpdateRackLevelDTO) (*entities.WarehouseRackLevel, error)
	DeleteRackLevel(ctx context.Context, businessID, levelID uint) error

	GetWarehouseTree(ctx context.Context, businessID, warehouseID uint) (*dtos.WarehouseTreeDTO, error)
}

type UseCase struct {
	repo ports.IRepository
}

func New(repo ports.IRepository) IUseCase {
	return &UseCase{repo: repo}
}
