package response

type ItemStockResult struct {
	ProductID    string
	SKU          string
	Requested    int
	Processed    int
	Sufficient   bool
	ErrorMessage string
}

type OrderStockResult struct {
	OrderID     string
	BusinessID  uint
	WarehouseID uint
	Success     bool
	ItemResults []ItemStockResult
}
