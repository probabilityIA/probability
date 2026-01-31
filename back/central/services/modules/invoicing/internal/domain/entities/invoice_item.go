package entities

import "time"

// InvoiceItem representa un item/línea de una factura
// Entidad PURA de dominio - SIN TAGS de infraestructura
type InvoiceItem struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	InvoiceID uint
	ProductID *string // VARCHAR en BD

	// Información del item
	SKU         string
	Name        string
	Description *string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
	Currency    string

	// Impuestos y descuentos
	Tax      float64
	TaxRate  *float64
	Discount float64

	// Información adicional del proveedor
	ProviderItemID *string
	Metadata       map[string]interface{}
}
