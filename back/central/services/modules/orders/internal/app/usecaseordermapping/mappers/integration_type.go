package mappers

// GetIntegrationTypeID convierte el código de tipo de integración a ID numérico
func GetIntegrationTypeID(integrationType string) uint {
	switch integrationType {
	case "shopify":
		return 1
	case "whatsapp", "whatsap", "whastap":
		return 2
	case "mercado_libre", "mercadolibre":
		return 3
	case "woocommerce", "woocormerce":
		return 4
	default:
		return 0
	}
}
