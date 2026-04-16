package dtos

// CreatePaymentDTO contiene los datos necesarios para iniciar un pago
type CreatePaymentDTO struct {
	BusinessID    uint
	Amount        float64
	Currency      string
	GatewayCode   string
	PaymentMethod string
	Description   string
	CallbackURL   *string
	Metadata      map[string]interface{}
}
