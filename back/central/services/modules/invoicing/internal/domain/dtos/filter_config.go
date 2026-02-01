package dtos

import "time"

// FilterConfig representa los filtros de configuración de facturación
// Se usa para deserializar el JSON de filtros
// NOTA: Este DTO es idéntico a entities.FilterConfig, se mantiene por compatibilidad
type FilterConfig struct {
	// Monto
	MinAmount *float64 `json:"min_amount,omitempty"`
	MaxAmount *float64 `json:"max_amount,omitempty"`

	// Pago
	PaymentStatus  *string `json:"payment_status,omitempty"`  // "paid", "unpaid", "partial"
	PaymentMethods []uint  `json:"payment_methods,omitempty"` // IDs de métodos permitidos

	// Orden
	OrderTypes      []string `json:"order_types,omitempty"`      // ["delivery", "pickup"]
	ExcludeStatuses []string `json:"exclude_statuses,omitempty"` // ["cancelled", "refunded"]

	// Productos
	ExcludeProducts     []string `json:"exclude_products,omitempty"`      // SKUs a excluir
	IncludeProductsOnly []string `json:"include_products_only,omitempty"` // Solo estos SKUs
	MinItemsCount       *int     `json:"min_items_count,omitempty"`
	MaxItemsCount       *int     `json:"max_items_count,omitempty"`

	// Cliente
	CustomerTypes      []string `json:"customer_types,omitempty"`       // ["natural", "juridica"]
	ExcludeCustomerIDs []string `json:"exclude_customer_ids,omitempty"` // IDs de clientes a excluir

	// Ubicación
	ShippingRegions []string `json:"shipping_regions,omitempty"` // ["Bogotá", "Medellín"]

	// Fecha
	DateRange *DateRangeFilter `json:"date_range,omitempty"`
}

type DateRangeFilter struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
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
