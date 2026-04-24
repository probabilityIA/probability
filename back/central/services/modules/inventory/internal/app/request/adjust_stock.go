package request

type AdjustStockDTO struct {
	ProductID   string
	WarehouseID uint
	LocationID  *uint
	LotID       *uint
	StateID     *uint
	UomID       *uint
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
	LotID           *uint
	StateID         *uint
	UomID           *uint
	BusinessID      uint
	Quantity        int
	Reason          string
	Notes           string
	CreatedByID     *uint
}
