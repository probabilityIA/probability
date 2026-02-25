package dtos

// CompareRequestDTO datos para solicitar comparación de facturas contra el proveedor
type CompareRequestDTO struct {
	DateFrom   string // YYYY-MM-DD
	DateTo     string // YYYY-MM-DD (máx 7 días después de DateFrom)
	BusinessID uint
}

// CompareItemDetail ítem de factura para comparación (sistema o proveedor)
type CompareItemDetail struct {
	ItemCode  string
	ItemName  string
	Quantity  string
	UnitValue string
	IVA       string
}

// CompareResult resultado de comparación para una factura
type CompareResult struct {
	Status          string   // "matched" | "system_only" | "provider_only"
	InvoiceNumber   string
	Prefix          string   // prefijo del documento (ej: "FEV")
	DocumentDate    string
	ProviderTotal   string   // string (Softpymes retorna string)
	SystemInvoiceID *uint
	SystemOrderID   *string
	SystemTotal     *float64
	CustomerNit     string
	CustomerName    string   // nombre del cliente del proveedor
	Comment         string   // "order:xxx" del campo comment de Softpymes
	OrderCreatedAt  *string              // fecha creación de la orden (YYYY-MM-DD)
	ProviderDetails []CompareItemDetail  // ítems de Softpymes
	SystemItems     []CompareItemDetail  // ítems de la factura en sistema
}

// CompareSummary resumen de la comparación
type CompareSummary struct {
	Matched      int
	SystemOnly   int
	ProviderOnly int
}

// CompareResponseData datos completos de la comparación (publicado por SSE)
type CompareResponseData struct {
	CorrelationID string
	BusinessID    uint
	DateFrom      string
	DateTo        string
	Results       []CompareResult
	Summary       CompareSummary
}

// Compare statuses
const (
	CompareStatusMatched      = "matched"
	CompareStatusSystemOnly   = "system_only"
	CompareStatusProviderOnly = "provider_only"
)
