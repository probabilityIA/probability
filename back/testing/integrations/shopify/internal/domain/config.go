package domain

// BusinessConfig contiene la configuración del business de prueba y su integración
type BusinessConfig struct {
	// Business data
	BusinessID   uint   `json:"business_id"`
	BusinessName string `json:"business_name"`
	BusinessCode string `json:"business_code"`

	// Integration data
	IntegrationID     uint   `json:"integration_id"`
	IntegrationName   string `json:"integration_name"`
	IntegrationCode   string `json:"integration_code"`
	IntegrationTypeID uint   `json:"integration_type_id"`

	// Shopify specific
	ShopDomain string `json:"shop_domain"`
	StoreID    string `json:"store_id"`
	APIVersion string `json:"api_version"`

	// Currency config
	ShopCurrency        string  `json:"shop_currency"`        // Moneda de la tienda Shopify (ej: "USD", "COP")
	PresentmentCurrency string  `json:"presentment_currency"` // Moneda de presentación al comprador (ej: "COP")
	TaxesIncluded       bool    `json:"taxes_included"`       // true si los precios incluyen IVA
	ExchangeRate        float64 `json:"exchange_rate"`        // Tasa de cambio PresentmentCurrency/ShopCurrency
}

// IsDualCurrency retorna true si la tienda usa una moneda diferente a la del comprador
func (c *BusinessConfig) IsDualCurrency() bool {
	return c.ShopCurrency != "" && c.PresentmentCurrency != "" && c.ShopCurrency != c.PresentmentCurrency
}

// DefaultTestBusinessConfig retorna la configuración del business de prueba (single-currency COP)
// Esta configuración debe coincidir con los datos en la base de datos
func DefaultTestBusinessConfig() *BusinessConfig {
	return &BusinessConfig{
		// Business: probability-dev (ID=7)
		BusinessID:   7,
		BusinessName: "probability-dev",
		BusinessCode: "probability-dev",

		// Integration: Shopify - pruebas (ID=1)
		IntegrationID:     1,
		IntegrationName:   "Shopify - pruebas",
		IntegrationCode:   "Shopify - pruebas",
		IntegrationTypeID: 1,

		// Shopify config
		ShopDomain: "tienda-test.myshopify.com",
		StoreID:    "tienda-test.myshopify.com",
		APIVersion: "2024-01",

		// Single currency (default, sin cambios al flujo actual)
		ShopCurrency:        "COP",
		PresentmentCurrency: "COP",
		TaxesIncluded:       false,
		ExchangeRate:        1.0,
	}
}

// DualCurrencyTestBusinessConfig retorna la configuración dual-currency USD/COP
// Simula una tienda Shopify configurada en USD con compradores colombianos (COP)
func DualCurrencyTestBusinessConfig() *BusinessConfig {
	return &BusinessConfig{
		// Business: probability-dev (ID=7) — mismo business de testing
		BusinessID:   7,
		BusinessName: "probability-dev",
		BusinessCode: "probability-dev",

		// Integration: Shopify - pruebas (ID=1)
		IntegrationID:     1,
		IntegrationName:   "Shopify - pruebas",
		IntegrationCode:   "Shopify - pruebas",
		IntegrationTypeID: 1,

		// Shopify config
		ShopDomain: "tienda-test.myshopify.com",
		StoreID:    "tienda-test.myshopify.com",
		APIVersion: "2024-01",

		// Dual currency: tienda en USD, compradores pagan en COP
		ShopCurrency:        "USD",
		PresentmentCurrency: "COP",
		TaxesIncluded:       true,
		ExchangeRate:        4200.0,
	}
}
