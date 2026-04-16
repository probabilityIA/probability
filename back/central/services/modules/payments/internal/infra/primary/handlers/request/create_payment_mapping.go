package request

// CreatePaymentMapping representa la solicitud HTTP para crear un mapeo
type CreatePaymentMapping struct {
	IntegrationType string `json:"integration_type" binding:"required,oneof=shopify whatsapp mercadolibre"`
	OriginalMethod  string `json:"original_method" binding:"required,max=128"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	Priority        int    `json:"priority"`
}
