package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	AdjustStock(ctx context.Context, dto request.AdjustStockDTO) (*entities.StockMovement, error)
	TransferStock(ctx context.Context, dto request.TransferStockDTO) error
	ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)

	ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	CreateMovementType(ctx context.Context, dto request.CreateStockMovementTypeDTO) (*entities.StockMovementType, error)
	UpdateMovementType(ctx context.Context, dto request.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error)
	DeleteMovementType(ctx context.Context, id uint) error

	BulkLoadInventory(ctx context.Context, dto request.BulkLoadDTO) (*response.BulkLoadResult, error)

	ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ConfirmSaleForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)

	ValidateCubing(ctx context.Context, dto request.ValidateCubingDTO) (*response.CubingCheckResult, error)
}

type useCase struct {
	repo           ports.IRepository
	publisher      ports.ISyncPublisher
	eventPublisher ports.IInventoryEventPublisher
	log            log.ILogger
}

func New(repo ports.IRepository, publisher ports.ISyncPublisher, eventPublisher ports.IInventoryEventPublisher, logger log.ILogger) IUseCase {
	return &useCase{
		repo:           repo,
		publisher:      publisher,
		eventPublisher: eventPublisher,
		log:            logger,
	}
}
