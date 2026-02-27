package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define los casos de uso del módulo inventory
type IUseCase interface {
	GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	AdjustStock(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error)
	TransferStock(ctx context.Context, dto dtos.TransferStockDTO) error
	ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)

	// Stock Movement Types
	ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	CreateMovementType(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error)
	UpdateMovementType(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error)
	DeleteMovementType(ctx context.Context, id uint) error

	// Order-driven inventory operations
	ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ConfirmSaleForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
	ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error)
}

// UseCase implementa IUseCase
type UseCase struct {
	repo           ports.IRepository
	publisher      ports.ISyncPublisher
	eventPublisher ports.IInventoryEventPublisher
	log            log.ILogger
}

// New crea una nueva instancia del use case
func New(repo ports.IRepository, publisher ports.ISyncPublisher, eventPublisher ports.IInventoryEventPublisher, logger log.ILogger) IUseCase {
	return &UseCase{
		repo:           repo,
		publisher:      publisher,
		eventPublisher: eventPublisher,
		log:            logger,
	}
}
