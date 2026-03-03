package dtos

// InvoiceConfigData representa la configuración adicional de facturación
type InvoiceConfigData struct {
	// Si incluye costo de envío
	IncludeShipping *bool `json:"include_shipping,omitempty"`

	// Si aplica descuentos
	ApplyDiscount *bool `json:"apply_discount,omitempty"`

	// Tasa de impuesto por defecto
	DefaultTaxRate *float64 `json:"default_tax_rate,omitempty"`

	// Tipo de factura
	InvoiceType *string `json:"invoice_type,omitempty"`

	// Notas por defecto
	Notes *string `json:"notes,omitempty"`
}
