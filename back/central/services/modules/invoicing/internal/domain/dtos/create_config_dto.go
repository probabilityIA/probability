package dtos

// CreateConfigDTO contiene los datos para crear una configuración de facturación
type CreateConfigDTO struct {
	// ID del negocio
	BusinessID uint

	// ID de la integración (fuente de órdenes)
	IntegrationID uint

	// ID del proveedor de facturación
	InvoicingProviderID uint

	// Si la configuración está habilitada
	Enabled bool

	// Si debe facturar automáticamente
	AutoInvoice bool

	// Filtros (ej: monto mínimo, métodos de pago permitidos)
	Filters map[string]interface{}

	// Configuración adicional de facturación
	InvoiceConfig map[string]interface{}

	// Descripción (opcional)
	Description *string

	// ID del usuario que crea la configuración
	CreatedByUserID uint
}
