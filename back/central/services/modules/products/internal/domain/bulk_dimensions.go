package domain

type DimensionRow struct {
	SKU    string
	Weight *float64
	Length *float64
	Width  *float64
	Height *float64
}

type BulkDimensionsResult struct {
	TotalRows    int      `json:"total_rows"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors,omitempty"`
}
