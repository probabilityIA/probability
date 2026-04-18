package request

type AdjustStockDTO struct {
	ProductID   string
	WarehouseID uint
	LocationID  *uint
	BusinessID  uint
	Quantity    int
	Reason      string
	Notes       string
	CreatedByID *uint
}

type TransferStockDTO struct {
	ProductID       string
	FromWarehouseID uint
	ToWarehouseID   uint
	FromLocationID  *uint
	ToLocationID    *uint
	BusinessID      uint
	Quantity        int
	Reason          string
	Notes           string
	CreatedByID     *uint
}
