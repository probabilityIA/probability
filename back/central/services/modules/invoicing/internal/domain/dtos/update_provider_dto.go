package dtos

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
