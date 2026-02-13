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
	InvoicingIntegrationID *uint // Integración de facturación (Softpymes - desde integrations/)

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
	ProviderImageURL *string // URL del logo del proveedor de facturación
}

// FilterConfig es la configuración completa de filtros
// NOTA: Réplica simplificada de la estructura del módulo de invoicing
type FilterConfig struct {
	// Monto
	MinAmount float64 `json:"min_amount,omitempty"`
	MaxAmount float64 `json:"max_amount,omitempty"`

	// Pago
	PaymentStatus  string `json:"payment_status,omitempty"`  // "paid", "unpaid", "partial"
	PaymentMethods []uint `json:"payment_methods,omitempty"` // IDs de métodos permitidos

	// Orden
	OrderTypes      []string `json:"order_types,omitempty"`      // ["delivery", "pickup"]
	ExcludeStatuses []string `json:"exclude_statuses,omitempty"` // ["cancelled", "refunded"]

	// Productos
	ExcludeProducts     []string `json:"exclude_products,omitempty"`      // SKUs a excluir
	IncludeProductsOnly []string `json:"include_products_only,omitempty"` // Solo estos SKUs
	MinItemsCount       int      `json:"min_items_count,omitempty"`
	MaxItemsCount       int      `json:"max_items_count,omitempty"`

	// Cliente
	CustomerTypes      []string `json:"customer_types,omitempty"`       // ["natural", "juridica"]
	ExcludeCustomerIDs []string `json:"exclude_customer_ids,omitempty"` // IDs de clientes a excluir

	// Ubicación
	ShippingRegions []string `json:"shipping_regions,omitempty"` // ["Bogotá", "Medellín"]
}
