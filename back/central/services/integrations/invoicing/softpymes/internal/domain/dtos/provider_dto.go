package dtos

// CreateProviderDTO contiene los datos para crear un proveedor de facturación
type CreateProviderDTO struct {
	// ID del negocio
	BusinessID uint

	// Código del tipo de proveedor (ej: "softpymes", "siigo")
	ProviderTypeCode string

	// Nombre del proveedor (ej: "Softpymes - Tienda Principal")
	Name string

	// Descripción (opcional)
	Description *string

	// Configuración específica del proveedor (no encriptada)
	Config map[string]interface{}

	// Credenciales (serán encriptadas)
	Credentials map[string]interface{}

	// Si es el proveedor por defecto
	IsDefault bool

	// ID del usuario que crea el proveedor
	CreatedByUserID uint
}

// UpdateProviderDTO contiene los datos para actualizar un proveedor
type UpdateProviderDTO struct {
	// Nombre del proveedor
	Name *string

	// Descripción
	Description *string

	// Configuración específica del proveedor
	Config map[string]interface{}

	// Credenciales (serán encriptadas)
	Credentials map[string]interface{}

	// Si está activo
	IsActive *bool

	// Si es el proveedor por defecto
	IsDefault *bool
}

// ProviderFiltersDTO contiene filtros para listar proveedores
type ProviderFiltersDTO struct {
	BusinessID       *uint
	ProviderTypeCode *string
	IsActive         *bool
	IsDefault        *bool
	Limit            int
	Offset           int
}
