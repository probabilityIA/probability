package domain

import "strings"

const (
	IntegrationTypeShopify      = 1
	IntegrationTypeWhatsApp     = 2
	IntegrationTypeMercadoLibre = 3
	IntegrationTypeWoocommerce  = 4
	IntegrationTypeInvoicing    = 5
	IntegrationTypePlatform     = 6
	IntegrationTypeFactus       = 7
	IntegrationTypeSiigo        = 8
	IntegrationTypeAlegra       = 9
	IntegrationTypeWorldOffice  = 10
	IntegrationTypeHelisa       = 11
	IntegrationTypeEnvioClick   = 12
	IntegrationTypeEnviame      = 13
	IntegrationTypeTu           = 14
	IntegrationTypeMiPaquete    = 15
	IntegrationTypeVTEX         = 16
	IntegrationTypeTiendanube   = 17
	IntegrationTypeMagento      = 18
	IntegrationTypeAmazon       = 19
	IntegrationTypeFalabella    = 20
	IntegrationTypeExito        = 21
	IntegrationTypeEmail        = 29
	IntegrationTypeTienda       = 30
	IntegrationTypeTiendaWeb    = 31
	IntegrationTypeJumpseller   = 33
)

const (
	IntegrationCategoryInternal = "internal"
	IntegrationCategoryExternal = "external"
)

func IntegrationTypeCodeAsInt(code string) int {
	lowerCode := strings.ToLower(code)
	switch lowerCode {
	case "shopify":
		return IntegrationTypeShopify
	case "whatsapp", "whatsap", "whastap":
		return IntegrationTypeWhatsApp
	case "mercado_libre", "mercado libre":
		return IntegrationTypeMercadoLibre
	case "woocommerce", "woocormerce":
		return IntegrationTypeWoocommerce
	case "softpymes", "invoicing":
		return IntegrationTypeInvoicing
	case "platform":
		return IntegrationTypePlatform
	case "factus":
		return IntegrationTypeFactus
	case "siigo":
		return IntegrationTypeSiigo
	case "alegra":
		return IntegrationTypeAlegra
	case "world_office", "worldoffice":
		return IntegrationTypeWorldOffice
	case "helisa":
		return IntegrationTypeHelisa
	case "envioclick":
		return IntegrationTypeEnvioClick
	case "enviame":
		return IntegrationTypeEnviame
	case "tu":
		return IntegrationTypeTu
	case "mipaquete", "mi_paquete":
		return IntegrationTypeMiPaquete
	case "vtex":
		return IntegrationTypeVTEX
	case "tiendanube", "tienda_nube":
		return IntegrationTypeTiendanube
	case "magento", "adobe_commerce":
		return IntegrationTypeMagento
	case "amazon":
		return IntegrationTypeAmazon
	case "falabella":
		return IntegrationTypeFalabella
	case "exito":
		return IntegrationTypeExito
	case "jumpseller":
		return IntegrationTypeJumpseller
	case "email":
		return IntegrationTypeEmail
	case "tienda":
		return IntegrationTypeTienda
	case "tienda_web":
		return IntegrationTypeTiendaWeb
	default:
		return 0
	}
}

func IsValidCategory(category string) bool {
	return category == IntegrationCategoryInternal || category == IntegrationCategoryExternal
}
