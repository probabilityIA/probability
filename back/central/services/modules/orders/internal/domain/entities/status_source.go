package entities

import "strings"

const (
	StatusSourceUser         = "user"
	StatusSourceSalesChannel = "sales_channel"
	StatusSourceCarrier      = "carrier"
	StatusSourceInventory    = "inventory"
	StatusSourceSystem       = "system"
)

const carrierChangedByLabel = "Transportadora"

var salesChannelLabels = map[string]string{
	"woocommerce":  "WooCommerce",
	"shopify":      "Shopify",
	"mercadolibre": "Mercado Libre",
	"meli":         "Mercado Libre",
	"jumpseller":   "Jumpseller",
	"vtex":         "VTEX",
	"tiendanube":   "Tiendanube",
	"magento":      "Magento",
	"amazon":       "Amazon",
	"falabella":    "Falabella",
	"exito":        "Exito",
}

// CarrierChangedBy es la etiqueta para cambios que llegan del tracking.
// Nunca se expone el nombre del proveedor de transporte.
func CarrierChangedBy() string {
	return carrierChangedByLabel
}

// SalesChannelChangedBy traduce el codigo de plataforma al nombre del canal.
func SalesChannelChangedBy(platform string) string {
	key := strings.ToLower(strings.TrimSpace(platform))
	if label, ok := salesChannelLabels[key]; ok {
		return label
	}
	if key == "" {
		return "Canal de venta"
	}
	return platform
}
