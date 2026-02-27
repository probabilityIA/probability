package dtos

// OrderInventoryItem representa un item de una orden para operaciones de inventario
type OrderInventoryItem struct {
	ProductID string
	SKU       string
	Quantity  int
}

// ReserveStockTxParams parámetros para la transacción de reserva de stock
type ReserveStockTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

// ReserveStockTxResult resultado de la transacción de reserva
type ReserveStockTxResult struct {
	PreviousAvailable int
	NewAvailable      int
	NewReserved       int
	Reserved          int  // cantidad efectivamente reservada
	Sufficient        bool // si había stock suficiente para reservar todo
}

// ConfirmSaleTxParams parámetros para confirmar venta (shipped/completed)
type ConfirmSaleTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

// ReleaseTxParams parámetros para liberar reserva (cancelled)
type ReleaseTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

// ReturnStockTxParams parámetros para devolver stock (refunded)
type ReturnStockTxParams struct {
	ProductID      string
	WarehouseID    uint
	BusinessID     uint
	Quantity       int
	MovementTypeID uint
	OrderID        string
}

// OrderStockResult resultado consolidado de una operación de inventario por orden
type OrderStockResult struct {
	OrderID     string
	BusinessID  uint
	WarehouseID uint
	Success     bool
	ItemResults []ItemStockResult
}

// ItemStockResult resultado de la operación de inventario para un item
type ItemStockResult struct {
	ProductID    string
	SKU          string
	Requested    int
	Processed    int
	Sufficient   bool
	ErrorMessage string
}
