package ports

import (
	"context"
)

// NequiConfig contiene las credenciales de Nequi
type NequiConfig struct {
	APIKey      string
	Environment string // "sandbox" | "production"
	PhoneCode   string // ej: "NIT_1"
}

// INequiClient define las operaciones del cliente HTTP de Nequi
type INequiClient interface {
	GenerateQR(ctx context.Context, config *NequiConfig, amount float64, reference string) (qrValue string, transactionID string, err error)
}

// IIntegrationRepository obtiene credenciales de Nequi desde integration_types
type IIntegrationRepository interface {
	GetNequiConfig(ctx context.Context) (*NequiConfig, error)
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
