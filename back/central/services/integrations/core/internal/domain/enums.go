package domain

// IntegrationType representa los tipos de integración disponibles
const (
	IntegrationTypeShopify      = 1
	IntegrationTypeWhatsApp     = 2
	IntegrationTypeMercadoLibre = 3
)

// IntegrationCategory representa la categoría de integración
const (
	IntegrationCategoryInternal = "internal" // Integraciones internas del sistema
	IntegrationCategoryExternal = "external" // Integraciones externas con clientes
)

// IsValidType valida si un tipo de integración es válido
func IsValidType(integrationType int) bool {
	validTypes := []int{
		IntegrationTypeShopify,
		IntegrationTypeWhatsApp,
		IntegrationTypeMercadoLibre,
	}
	for _, validType := range validTypes {
		if integrationType == validType {
			return true
		}
	}
	return false
}

// IsValidCategory valida si una categoría de integración es válida
func IsValidCategory(category string) bool {
	return category == IntegrationCategoryInternal || category == IntegrationCategoryExternal
}
