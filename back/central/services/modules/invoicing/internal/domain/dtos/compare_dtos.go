package dtos

// CompareRequestDTO datos para solicitar comparación de facturas contra el proveedor
type CompareRequestDTO struct {
	DateFrom   string // YYYY-MM-DD
	DateTo     string // YYYY-MM-DD (máx 7 días después de DateFrom)
	BusinessID uint
}

// CompareItemDetail ítem de factura para comparación (sistema o proveedor)
type CompareItemDetail struct {
	ItemCode  string `json:"item_code"`
	ItemName  string `json:"item_name"`
	Quantity  string `json:"quantity"`
	UnitValue string `json:"unit_value"`
	IVA       string `json:"iva"`
}

// CompareResult resultado de comparación para una factura
type CompareResult struct {
	Status          string              `json:"status"`           // "matched" | "system_only" | "provider_only"
	InvoiceNumber   string              `json:"invoice_number"`
	Prefix          string              `json:"prefix"`           // prefijo del documento (ej: "FEV")
	DocumentDate    string              `json:"document_date"`
	ProviderTotal   string              `json:"provider_total"`   // string (Softpymes retorna string)
	SystemInvoiceID *uint               `json:"system_invoice_id"`
	SystemOrderID   *string             `json:"system_order_id"`
	SystemTotal     *float64            `json:"system_total"`
	CustomerNit     string              `json:"customer_nit"`
	CustomerName    string              `json:"customer_name"`    // nombre del cliente del proveedor
	Comment         string              `json:"comment"`          // "order:xxx" del campo comment de Softpymes
	OrderCreatedAt  *string             `json:"order_created_at"` // fecha creación de la orden (YYYY-MM-DD)
	ProviderDetails []CompareItemDetail `json:"provider_details"` // ítems de Softpymes
	SystemItems     []CompareItemDetail `json:"system_items"`     // ítems de la factura en sistema
}

// CompareSummary resumen de la comparación
type CompareSummary struct {
	Matched      int `json:"matched"`
	SystemOnly   int `json:"system_only"`
	ProviderOnly int `json:"provider_only"`
}

// CompareResponseData datos completos de la comparación (publicado por SSE)
type CompareResponseData struct {
	CorrelationID string          `json:"correlation_id"`
	BusinessID    uint            `json:"business_id"`
	DateFrom      string          `json:"date_from"`
	DateTo        string          `json:"date_to"`
	Results       []CompareResult `json:"results"`
	Summary       CompareSummary  `json:"summary"`
}

// Compare statuses
const (
	CompareStatusMatched      = "matched"
	CompareStatusSystemOnly   = "system_only"
	CompareStatusProviderOnly = "provider_only"
)
