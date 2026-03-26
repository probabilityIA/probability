package entities

import "time"

// Invoice representa una factura electrónica generada para una orden
// Entidad PURA de dominio - SIN TAGS de infraestructura
type Invoice struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	OrderID                string // UUID de la orden
	BusinessID             uint
	InvoicingProviderID    *uint  // DEPRECATED: Mantener temporalmente para dual-read
	InvoicingIntegrationID *uint  // Integración de facturación (Softpymes - desde integrations/)

	// Identificadores
	InvoiceNumber  string
	ExternalID     *string
	InternalNumber string

	// Información financiera
	Subtotal     float64
	Tax          float64
	Discount     float64
	ShippingCost     float64
	ShippingDiscount float64
	TotalAmount  float64
	Currency     string

	// Información del cliente
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	CustomerDNI   string

	// Estado
	Status string

	// Testing
	IsTest bool

	// Timestamps
	IssuedAt    *time.Time
	CancelledAt *time.Time
	ExpiresAt   *time.Time

	// URLs y archivos
	InvoiceURL *string
	PDFURL     *string
	XMLURL     *string
	CUFE       *string

	// Información adicional
	Notes            *string
	Metadata         map[string]interface{}
	ProviderResponse map[string]interface{}

	// Datos derivados de relaciones (no almacenados en tabla invoices)
	OrderNumber     string  // Número visible de la orden (de tabla orders)
	ProviderLogoURL *string // Logo del proveedor de facturación
	ProviderName    *string // Nombre del tipo de integración de facturación

	// Items de la factura (relación)
	Items []InvoiceItem
}
