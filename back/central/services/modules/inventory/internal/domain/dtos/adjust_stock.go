package dtos

// AdjustStockDTO datos para ajustar stock manualmente
type AdjustStockDTO struct {
	ProductID   string
	WarehouseID uint
	LocationID  *uint
	BusinessID  uint
	Quantity    int    // positivo=agregar, negativo=quitar
	Reason      string // motivo del ajuste
	Notes       string
	CreatedByID *uint
}

// TransferStockDTO datos para transferir stock entre bodegas
type TransferStockDTO struct {
	ProductID         string
	FromWarehouseID   uint
	ToWarehouseID     uint
	FromLocationID    *uint
	ToLocationID      *uint
	BusinessID        uint
	Quantity          int // siempre positivo
	Reason            string
	Notes             string
	CreatedByID       *uint
}
