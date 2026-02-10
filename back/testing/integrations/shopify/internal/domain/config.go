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
}

// DefaultTestBusinessConfig retorna la configuración del business de prueba
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
	}
}
