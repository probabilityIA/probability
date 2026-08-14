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

type FamilyImportMatch struct {
	ProductID         string   `json:"product_id"`
	SKU               string   `json:"sku"`
	Name              string   `json:"name"`
	CurrentFamilyID   *uint    `json:"current_family_id,omitempty"`
	CurrentFamilyName string   `json:"current_family_name,omitempty"`
	Weight            *float64 `json:"weight,omitempty"`
	Length            *float64 `json:"length,omitempty"`
	Width             *float64 `json:"width,omitempty"`
	Height            *float64 `json:"height,omitempty"`
}

type FamilyImportPreview struct {
	Matched  []FamilyImportMatch `json:"matched"`
	NotFound []string            `json:"not_found"`
}
