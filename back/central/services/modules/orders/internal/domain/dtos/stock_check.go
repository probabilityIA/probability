package dtos

type StockCheckItem struct {
	ProductID  string
	ProductSKU string
	Quantity   int
}

type StockCheckResult struct {
	ProductID  string
	ProductSKU string
	Required   int
	Available  int
	Sufficient bool
}
