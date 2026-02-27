package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// IRepository define los métodos del repositorio del módulo inventory
type IRepository interface {
	// Inventory Levels
	GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	GetOrCreateLevel(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error)
	UpdateLevel(ctx context.Context, level *entities.InventoryLevel) error

	// Stock Movements
	CreateMovement(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error)
	ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)

	// Stock Movement Types
	ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	GetMovementTypeByID(ctx context.Context, id uint) (*entities.StockMovementType, error)
	GetMovementTypeIDByCode(ctx context.Context, code string) (uint, error)
	CreateMovementType(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error)
	UpdateMovementType(ctx context.Context, movType *entities.StockMovementType) error
	DeleteMovementType(ctx context.Context, id uint) error

	// Product queries (replicadas localmente)
	GetProductByID(ctx context.Context, productID string, businessID uint) (productName string, productSKU string, trackInventory bool, err error)
	UpdateProductStockQuantity(ctx context.Context, productID string, totalQuantity int) error

	// Warehouse validation (replicada localmente)
	WarehouseExists(ctx context.Context, warehouseID uint, businessID uint) (bool, error)

	// Product-Integration queries (replicada para sync)
	GetProductIntegrations(ctx context.Context, productID string, businessID uint) ([]ProductIntegrationInfo, error)

	// Operaciones transaccionales (SELECT FOR UPDATE + commit atómico)
	AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error)
	TransferStockTx(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error)

	// Operaciones transaccionales de órdenes
	ReserveStockTx(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error)
	ConfirmSaleTx(ctx context.Context, params dtos.ConfirmSaleTxParams) error
	ReleaseStockTx(ctx context.Context, params dtos.ReleaseTxParams) error
	ReturnStockTx(ctx context.Context, params dtos.ReturnStockTxParams) error

	// Warehouse queries
	GetDefaultWarehouseID(ctx context.Context, businessID uint) (uint, error)
}

// ProductIntegrationInfo datos de una integración vinculada a un producto
type ProductIntegrationInfo struct {
	IntegrationID       uint
	ExternalProductID   string
	IntegrationTypeCode string
}

// ISyncPublisher publica mensajes de sync a RabbitMQ
type ISyncPublisher interface {
	PublishInventorySync(ctx context.Context, msg InventorySyncMessage) error
}

// IInventoryEventPublisher publica eventos de inventario a Redis SSE
type IInventoryEventPublisher interface {
	PublishInventoryEvent(ctx context.Context, event InventoryEvent) error
}

// InventoryEvent evento de inventario para SSE
type InventoryEvent struct {
	EventType   string                 `json:"event_type"`
	OrderID     string                 `json:"order_id"`
	BusinessID  uint                   `json:"business_id"`
	WarehouseID uint                   `json:"warehouse_id"`
	Timestamp   string                 `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// InventorySyncMessage mensaje que se publica a RabbitMQ para sincronizar inventario
type InventorySyncMessage struct {
	ProductID         string `json:"product_id"`
	ExternalProductID string `json:"external_product_id"`
	IntegrationID     uint   `json:"integration_id"`
	BusinessID        uint   `json:"business_id"`
	NewQuantity       int    `json:"new_quantity"`
	WarehouseID       uint   `json:"warehouse_id"`
	Source            string `json:"source"` // manual_adjustment, transfer, sync
	Timestamp         string `json:"timestamp"`
}
