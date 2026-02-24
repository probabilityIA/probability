package domain

import "strings"

// Integration type IDs — mapean a integration_types.id en la base de datos
const (
	IntegrationTypeShopify      = 1
	IntegrationTypeWhatsApp     = 2 // Whastap en la BD
	IntegrationTypeMercadoLibre = 3
	IntegrationTypeWoocommerce  = 4
	IntegrationTypeInvoicing    = 5 // Softpymes - Facturación electrónica
	IntegrationTypePlatform     = 6 // Plataforma interna
	IntegrationTypeFactus       = 7 // Factus - Facturación electrónica
	IntegrationTypeSiigo        = 8  // Siigo - Facturación electrónica
	IntegrationTypeAlegra       = 9  // Alegra - Facturación electrónica
	IntegrationTypeWorldOffice  = 10 // World Office - Facturación electrónica
	IntegrationTypeHelisa       = 11 // Helisa - Facturación electrónica
	IntegrationTypeEnvioClick   = 12 // EnvioClick - Transporte
	IntegrationTypeEnviame      = 13 // Enviame - Transporte
	IntegrationTypeTu           = 14 // Tu - Transporte
	IntegrationTypeMiPaquete    = 15 // MiPaquete - Transporte
	IntegrationTypeVTEX         = 16 // VTEX - E-commerce
	IntegrationTypeTiendanube   = 17 // Tiendanube - E-commerce
	IntegrationTypeMagento      = 18 // Magento/Adobe Commerce - E-commerce
	IntegrationTypeAmazon       = 19 // Amazon - Marketplace
	IntegrationTypeFalabella    = 20 // Falabella - Marketplace
	IntegrationTypeExito        = 21 // Exito - Marketplace
)

// IntegrationCategory representa la categoría de integración
const (
	IntegrationCategoryInternal = "internal" // Integraciones internas del sistema
	IntegrationCategoryExternal = "external" // Integraciones externas con clientes
)

// IntegrationTypeCodeAsInt convierte el código string de tipo de integración a su ID numérico.
// Esta es la función canónica — usarla en todo el proyecto.
func IntegrationTypeCodeAsInt(code string) int {
	lowerCode := strings.ToLower(code)
	switch lowerCode {
	case "shopify":
		return IntegrationTypeShopify // 1
	case "whatsapp", "whatsap", "whastap":
		return IntegrationTypeWhatsApp // 2
	case "mercado_libre", "mercado libre":
		return IntegrationTypeMercadoLibre // 3
	case "woocommerce", "woocormerce":
		return IntegrationTypeWoocommerce // 4
	case "softpymes", "invoicing":
		return IntegrationTypeInvoicing // 5
	case "platform":
		return IntegrationTypePlatform // 6
	case "factus":
		return IntegrationTypeFactus // 7
	case "siigo":
		return IntegrationTypeSiigo // 8
	case "alegra":
		return IntegrationTypeAlegra // 9
	case "world_office", "worldoffice":
		return IntegrationTypeWorldOffice // 10
	case "helisa":
		return IntegrationTypeHelisa // 11
	case "envioclick":
		return IntegrationTypeEnvioClick // 12
	case "enviame":
		return IntegrationTypeEnviame // 13
	case "tu":
		return IntegrationTypeTu // 14
	case "mipaquete", "mi_paquete":
		return IntegrationTypeMiPaquete // 15
	case "vtex":
		return IntegrationTypeVTEX // 16
	case "tiendanube", "tienda_nube":
		return IntegrationTypeTiendanube // 17
	case "magento", "adobe_commerce":
		return IntegrationTypeMagento // 18
	case "amazon":
		return IntegrationTypeAmazon // 19
	case "falabella":
		return IntegrationTypeFalabella // 20
	case "exito":
		return IntegrationTypeExito // 21
	default:
		return 0
	}
}

// IsValidCategory valida si una categoría de integración es válida
func IsValidCategory(category string) bool {
	return category == IntegrationCategoryInternal || category == IntegrationCategoryExternal
}
