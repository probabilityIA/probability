package request

type BulkLoadItem struct {
	SKU          string
	Quantity     int
	MinStock     *int
	MaxStock     *int
	ReorderPoint *int
}

type BulkLoadDTO struct {
	WarehouseID uint
	BusinessID  uint
	CreatedByID *uint
	Reason      string
	Items       []BulkLoadItem
}
