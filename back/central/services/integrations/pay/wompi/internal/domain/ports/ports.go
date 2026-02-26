package ports

import "context"

// WompiConfig contiene las credenciales de Wompi
type WompiConfig struct {
	PrivateKey  string
	Environment string // "sandbox" | "production"
}

// IWompiClient define las operaciones del cliente HTTP de Wompi
// TODO: ajustar método y parámetros según la API oficial de Wompi
// Docs: https://docs.wompi.co/
type IWompiClient interface {
	CreateTransaction(ctx context.Context, config *WompiConfig, amount float64, currency, reference, description string) (transactionID string, redirectURL string, err error)
}

// IIntegrationRepository obtiene credenciales de Wompi desde integration_types
type IIntegrationRepository interface {
	GetWompiConfig(ctx context.Context) (*WompiConfig, error)
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
