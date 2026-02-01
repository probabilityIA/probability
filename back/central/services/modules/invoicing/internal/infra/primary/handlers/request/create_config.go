package request

// CreateConfig es el request para crear una configuración de facturación
type CreateConfig struct {
	BusinessID             uint                   `json:"business_id" binding:"required"`
	IntegrationID          uint                   `json:"integration_id" binding:"required"`              // Shopify, MercadoLibre, etc.
	InvoicingIntegrationID uint                   `json:"invoicing_integration_id" binding:"required"`    // Integración que emitirá facturas (nueva columna)
	InvoicingProviderID    *uint                  `json:"invoicing_provider_id"`                          // Deprecado - mantener para compatibilidad
	Enabled                *bool                  `json:"enabled"`                                        // Por defecto true
	AutoInvoice            *bool                  `json:"auto_invoice"`                                   // Por defecto false (manual)
	Filters                map[string]interface{} `json:"filters"`                                        // min_amount, payment_status, etc.
}
