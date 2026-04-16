package ports

import "context"

// EPaycoConfig contiene las credenciales de ePayco
type EPaycoConfig struct {
	CustomerID  string // p_cust_id_cliente
	Key         string // p_key
	Environment string // "test" | "production"
}

// IEPaycoClient define las operaciones del cliente HTTP de ePayco
// TODO: ajustar método y parámetros según la API oficial de ePayco
// Docs: https://docs.epayco.co/
type IEPaycoClient interface {
	CreateCheckout(ctx context.Context, config *EPaycoConfig, amount float64, currency, reference, description string) (checkoutID string, redirectURL string, err error)
}

// IIntegrationRepository obtiene credenciales de ePayco desde integration_types
type IIntegrationRepository interface {
	GetEPaycoConfig(ctx context.Context) (*EPaycoConfig, error)
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
