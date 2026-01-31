package entities

import "time"

// InvoicingConfig define qué integraciones deben facturar automáticamente
// y con qué proveedor de facturación.
// Entidad PURA de dominio - SIN TAGS de infraestructura
type InvoicingConfig struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	BusinessID          uint
	IntegrationID       uint
	InvoicingProviderID uint

	// Estado
	Enabled     bool
	AutoInvoice bool

	// Filtros y configuración (serán JSON en infraestructura)
	Filters       map[string]interface{}
	InvoiceConfig map[string]interface{}

	// Metadata
	Description string
	CreatedByID uint
	UpdatedByID *uint
}
