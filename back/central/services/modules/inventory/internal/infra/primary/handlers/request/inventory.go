package request

// AdjustStockRequest payload para ajustar stock
type AdjustStockRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	WarehouseID uint   `json:"warehouse_id" binding:"required,min=1"`
	LocationID  *uint  `json:"location_id"`
	Quantity    int    `json:"quantity" binding:"required"`
	Reason      string `json:"reason" binding:"required,min=2,max=255"`
	Notes       string `json:"notes" binding:"omitempty,max=1000"`
}

// TransferStockRequest payload para transferir stock entre bodegas
type TransferStockRequest struct {
	ProductID       string `json:"product_id" binding:"required"`
	FromWarehouseID uint   `json:"from_warehouse_id" binding:"required,min=1"`
	ToWarehouseID   uint   `json:"to_warehouse_id" binding:"required,min=1"`
	FromLocationID  *uint  `json:"from_location_id"`
	ToLocationID    *uint  `json:"to_location_id"`
	Quantity        int    `json:"quantity" binding:"required,min=1"`
	Reason          string `json:"reason" binding:"omitempty,max=255"`
	Notes           string `json:"notes" binding:"omitempty,max=1000"`
}
