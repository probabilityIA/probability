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
	Status        string  `json:"status"`         // "matched", "provider_only", "system_only"
	ItemCode      string  `json:"item_code"`      // SKU / itemCode
	ProviderName  string  `json:"provider_name"`
	SystemName    string  `json:"system_name"`
	ProviderPrice float64 `json:"provider_price"`
	SystemPrice   float64 `json:"system_price"`
	PriceDiff     float64 `json:"price_diff"` // provider - system (0 si falta en alguno)
	UnitCost      float64 `json:"unit_cost"`
	Description   string  `json:"description"`
}

// ItemCompareSummary resumen de la comparación de ítems
type ItemCompareSummary struct {
	Matched       int `json:"matched"`
	ProviderOnly  int `json:"provider_only"`
	SystemOnly    int `json:"system_only"`
	TotalProvider int `json:"total_provider"`
	TotalSystem   int `json:"total_system"`
}

// ItemCompareResponseData datos completos de la comparación de ítems (publicado por SSE)
type ItemCompareResponseData struct {
	CorrelationID string             `json:"correlation_id"`
	BusinessID    uint               `json:"business_id"`
	Results       []ItemCompareResult `json:"results"`
	Summary       ItemCompareSummary  `json:"summary"`
}

// BANK ACCOUNTS

// ListBankAccountsRequestDTO datos para solicitar cuentas bancarias del proveedor
type ListBankAccountsRequestDTO struct {
	BusinessID uint
}

// BankAccountResult cuenta bancaria del proveedor
type BankAccountResult struct {
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	NameType      string `json:"name_type"`
}

// BankAccountsResponseData datos completos de las cuentas bancarias (almacenado en Redis)
type BankAccountsResponseData struct {
	CorrelationID string              `json:"correlation_id"`
	BusinessID    uint                `json:"business_id"`
	Results       []BankAccountResult `json:"results"`
}
