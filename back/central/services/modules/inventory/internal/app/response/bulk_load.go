package response

type BulkLoadItemResult struct {
	SKU         string
	ProductID   string
	Success     bool
	PreviousQty int
	NewQty      int
	Error       string
}

type BulkLoadResult struct {
	TotalItems   int
	SuccessCount int
	FailureCount int
	Items        []BulkLoadItemResult
}
