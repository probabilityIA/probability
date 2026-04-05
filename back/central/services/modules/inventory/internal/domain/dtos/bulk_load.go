package dtos

// BulkLoadItem un item individual de la carga masiva
type BulkLoadItem struct {
	SKU          string
	Quantity     int
	MinStock     *int
	MaxStock     *int
	ReorderPoint *int
}

// BulkLoadDTO datos para carga masiva de inventario
type BulkLoadDTO struct {
	WarehouseID uint
	BusinessID  uint
	CreatedByID *uint
	Reason      string
	Items       []BulkLoadItem
}

// BulkLoadResult resultado de la carga masiva
type BulkLoadResult struct {
	TotalItems   int
	SuccessCount int
	FailureCount int
	Items        []BulkLoadItemResult
}

// BulkLoadItemResult resultado individual por item
type BulkLoadItemResult struct {
	SKU         string
	ProductID   string
	Success     bool
	PreviousQty int
	NewQty      int
	Error       string
}
