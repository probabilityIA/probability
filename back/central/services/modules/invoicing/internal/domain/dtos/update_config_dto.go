package dtos

// UpdateConfigDTO contiene los datos para actualizar una configuración
type UpdateConfigDTO struct {
	// Si la configuración está habilitada
	Enabled *bool

	// Si se debe facturar automáticamente
	AutoInvoice *bool

	// Filtros de configuración
	Filters map[string]interface{}
}
