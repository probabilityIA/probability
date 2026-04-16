package ports

import (
	"context"
)

// BoldConfig contiene las credenciales de Bold
type BoldConfig struct {
	APIKey      string
	Environment string // "sandbox" | "production"
}

// IBoldClient define las operaciones del cliente HTTP de Bold
type IBoldClient interface {
	CreatePaymentLink(ctx context.Context, config *BoldConfig, amount float64, currency, reference, description string) (linkID string, checkoutURL string, err error)
}

// IIntegrationRepository obtiene credenciales de Bold desde integration_types
type IIntegrationRepository interface {
	GetBoldConfig(ctx context.Context) (*BoldConfig, error)
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
