package request

// CreatePaymentRequest es el cuerpo de la petición para crear un pago
type CreatePaymentRequest struct {
	Amount        float64                `json:"amount" binding:"required,gt=0"`
	Currency      string                 `json:"currency"`
	GatewayCode   string                 `json:"gateway_code" binding:"required"`
	PaymentMethod string                 `json:"payment_method"`
	Description   string                 `json:"description"`
	CallbackURL   *string                `json:"callback_url"`
	Metadata      map[string]interface{} `json:"metadata"`
}
