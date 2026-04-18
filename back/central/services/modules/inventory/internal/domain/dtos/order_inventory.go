package dtos

type OrderInventoryItem struct {
	ProductID string
	SKU       string
	Quantity  int
}

type ReserveStockTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

type ReserveStockTxResult struct {
	PreviousAvailable int
	NewAvailable      int
	NewReserved       int
	Reserved          int
	Sufficient        bool
}

type ConfirmSaleTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

type ReleaseTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

type ReturnStockTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}
