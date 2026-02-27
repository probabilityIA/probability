package dtos

import "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"

// AdjustStockTxParams agrupa todo lo que necesita la transacción de ajuste
type AdjustStockTxParams struct {
	ProductID      string
	WarehouseID    uint
	LocationID     *uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	Reason         string
	Notes          string
	ReferenceType  string
	CreatedByID    *uint
}

// AdjustStockTxResult resultado de la transacción de ajuste
type AdjustStockTxResult struct {
	Movement    *entities.StockMovement
	NewQuantity int
	Level       *entities.InventoryLevel
}

// TransferStockTxParams agrupa todo lo que necesita la transacción de transferencia
type TransferStockTxParams struct {
	ProductID       string
	FromWarehouseID uint
	ToWarehouseID   uint
	FromLocationID  *uint
	ToLocationID    *uint
	BusinessID      uint
	Quantity        int
	MovementTypeID  uint
	Reason          string
	Notes           string
	ReferenceType   string
	CreatedByID     *uint
}

// TransferStockTxResult resultado de la transacción de transferencia
type TransferStockTxResult struct {
	FromNewQty int
	ToNewQty   int
}
