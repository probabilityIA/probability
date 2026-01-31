package dtos

// FilterConfig representa los filtros de configuración de facturación
// Se usa para deserializar el JSON de filtros
type FilterConfig struct {
	// Monto mínimo para facturar
	MinAmount *float64 `json:"min_amount,omitempty"`

	// Estado de pago requerido (ej: "paid")
	PaymentStatus *string `json:"payment_status,omitempty"`

	// IDs de métodos de pago permitidos
	PaymentMethods []uint `json:"payment_methods,omitempty"`

	// Tipos de orden permitidos (ej: ["delivery", "pickup"])
	OrderTypes []string `json:"order_types,omitempty"`

	// Estados de orden a excluir
	ExcludeStatuses []string `json:"exclude_statuses,omitempty"`
}

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
