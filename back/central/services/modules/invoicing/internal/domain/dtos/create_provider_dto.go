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
