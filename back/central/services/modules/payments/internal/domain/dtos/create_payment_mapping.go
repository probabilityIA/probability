package dtos

// CreatePaymentMapping representa la solicitud para crear un mapeo (PURO - sin tags)
type CreatePaymentMapping struct {
	IntegrationType string // shopify, whatsapp, mercadolibre
	OriginalMethod  string
	PaymentMethodID uint
	Priority        int
}
