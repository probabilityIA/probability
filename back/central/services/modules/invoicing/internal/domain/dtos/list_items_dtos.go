package dtos

// ListItemsRequestDTO datos para solicitar comparación de ítems del proveedor vs productos del sistema
type ListItemsRequestDTO struct {
	BusinessID uint
}

// ProviderItem ítem del catálogo del proveedor (Softpymes)
type ProviderItem struct {
	ItemCode      string
	ItemName      string
	ItemPrice     float64
	UnitCost      float64
	Description   string
	MinimumStock  string
	OrderQuantity string
}

// SystemProduct producto del sistema Probability
type SystemProduct struct {
	ID    string
	SKU   string
	Name  string
	Price float64
}

// ItemCompareResult resultado de comparación para un ítem/producto
type ItemCompareResult struct {
	Status        string  // "matched", "provider_only", "system_only"
	ItemCode      string  // SKU / itemCode
	ProviderName  string
	SystemName    string
	ProviderPrice float64
	SystemPrice   float64
	PriceDiff     float64 // provider - system (0 si falta en alguno)
	UnitCost      float64
	Description   string
}

// ItemCompareSummary resumen de la comparación de ítems
type ItemCompareSummary struct {
	Matched       int
	ProviderOnly  int
	SystemOnly    int
	TotalProvider int
	TotalSystem   int
}

// ItemCompareResponseData datos completos de la comparación de ítems (publicado por SSE)
type ItemCompareResponseData struct {
	CorrelationID string
	BusinessID    uint
	Results       []ItemCompareResult
	Summary       ItemCompareSummary
}
