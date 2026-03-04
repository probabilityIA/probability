package dtos

// UpdateConfigDTO contiene los datos para actualizar una configuración
type UpdateConfigDTO struct {
	// Si la configuración está habilitada
	Enabled *bool

	// Si se debe facturar automáticamente
	AutoInvoice *bool

	// ID de la integración de facturación (FK a integrations)
	InvoicingIntegrationID *uint

	// IDs de integraciones de e-commerce (fuentes de órdenes)
	// Si se provee, reemplaza todas las integraciones actuales de la config
	IntegrationIDs *[]uint

	// Filtros de configuración
	Filters map[string]interface{}

	// Business ID del solicitante (para validar pertenencia)
	// Super admin: proviene de ?business_id query param
	// Usuario normal: proviene del JWT
	RequestingBusinessID *uint
}
