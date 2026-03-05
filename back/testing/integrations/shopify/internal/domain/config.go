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

	// Currency config — simula tienda Shopify en USD con compradores colombianos (COP)
	ShopCurrency        string  `json:"shop_currency"`        // "USD"
	PresentmentCurrency string  `json:"presentment_currency"` // "COP"
	TaxesIncluded       bool    `json:"taxes_included"`       // true (precios incluyen IVA 19%)
	ExchangeRate        float64 `json:"exchange_rate"`        // ~4200 COP/USD
}

// IsDualCurrency retorna true si la tienda usa una moneda diferente a la del comprador
func (c *BusinessConfig) IsDualCurrency() bool {
	return c.ShopCurrency != "" && c.PresentmentCurrency != "" && c.ShopCurrency != c.PresentmentCurrency
}

// DefaultTestBusinessConfig retorna la configuración del business de prueba
// Simula una tienda Shopify configurada en USD con compradores colombianos (COP)
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

		// Dual currency: tienda en USD, compradores pagan en COP
		ShopCurrency:        "USD",
		PresentmentCurrency: "COP",
		TaxesIncluded:       true,
		ExchangeRate:        4200.0,
	}
}
