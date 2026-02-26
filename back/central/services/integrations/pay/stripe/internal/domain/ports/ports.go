package ports

import "context"

// StripeConfig contiene las credenciales de Stripe
type StripeConfig struct {
	SecretKey   string
	Environment string // "test" | "live"
}

// IStripeClient define las operaciones del cliente HTTP de Stripe
// TODO: ajustar método y parámetros según la API oficial de Stripe
// Docs: https://stripe.com/docs/api
type IStripeClient interface {
	CreatePaymentIntent(ctx context.Context, config *StripeConfig, amount float64, currency, reference, description string) (paymentIntentID string, clientSecret string, err error)
}

// IIntegrationRepository obtiene credenciales de Stripe desde integration_types
type IIntegrationRepository interface {
	GetStripeConfig(ctx context.Context) (*StripeConfig, error)
}

// IResponsePublisher publica respuestas al módulo de pagos
type IResponsePublisher interface {
	PublishPaymentResponse(ctx context.Context, msg *PaymentResponseMsg) error
}

// PaymentResponseMsg mensaje de respuesta a publicar en pay.responses
type PaymentResponseMsg struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	GatewayCode          string                 `json:"gateway_code"`
	Status               string                 `json:"status"` // "success"|"error"
	ExternalID           *string                `json:"external_id,omitempty"`
	GatewayResponse      map[string]interface{} `json:"gateway_response,omitempty"`
	Error                string                 `json:"error,omitempty"`
	ErrorCode            string                 `json:"error_code,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	ProcessingTimeMs     int64                  `json:"processing_time_ms"`
}
