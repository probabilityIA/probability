package entities

import "time"

// InvoicingConfig define qué integraciones deben facturar automáticamente
// y con qué proveedor de facturación.
// NOTA: Esta es una réplica local de la entidad del módulo de invoicing
// para evitar dependencias con paquetes internal/
type InvoicingConfig struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	BusinessID             uint
	IntegrationID          uint  // Integración de e-commerce (Shopify, MeLi, etc.)
	InvoicingProviderID    *uint // DEPRECATED: Mantener temporalmente para dual-read
	InvoicingIntegrationID *uint // Integración de facturación (World Office - desde integrations/)

	// Estado
	Enabled     bool
	AutoInvoice bool

	// Filtros y configuración
	Filters       FilterConfig
	InvoiceConfig map[string]interface{}

	// Metadata
	Description string
	CreatedByID uint
	UpdatedByID *uint

	// Nombres de relaciones (opcionales - populados por repo)
	IntegrationName  *string
	ProviderName     *string
	ProviderImageURL *string
}

// FilterConfig es la configuración completa de filtros
type FilterConfig struct {
	// Monto
	MinAmount float64 `json:"min_amount,omitempty"`
	MaxAmount float64 `json:"max_amount,omitempty"`

	// Pago
	PaymentStatus  string `json:"payment_status,omitempty"`
	PaymentMethods []uint `json:"payment_methods,omitempty"`

	// Orden
	OrderTypes      []string `json:"order_types,omitempty"`
	ExcludeStatuses []string `json:"exclude_statuses,omitempty"`

	// Productos
	ExcludeProducts     []string `json:"exclude_products,omitempty"`
	IncludeProductsOnly []string `json:"include_products_only,omitempty"`
	MinItemsCount       int      `json:"min_items_count,omitempty"`
	MaxItemsCount       int      `json:"max_items_count,omitempty"`

	// Cliente
	CustomerTypes      []string `json:"customer_types,omitempty"`
	ExcludeCustomerIDs []string `json:"exclude_customer_ids,omitempty"`

	// Ubicación
	ShippingRegions []string `json:"shipping_regions,omitempty"`
}
