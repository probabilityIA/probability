package ports

import "context"

// MeliPagoConfig contiene las credenciales de MercadoPago
type MeliPagoConfig struct {
	AccessToken string
	Environment string // "sandbox" | "production"
}

// IMeliPagoClient define las operaciones del cliente HTTP de MercadoPago
// TODO: ajustar método y parámetros según la API oficial de MercadoPago
// Docs: https://www.mercadopago.com.co/developers/es/reference
type IMeliPagoClient interface {
	CreatePreference(ctx context.Context, config *MeliPagoConfig, amount float64, currency, reference, description string) (preferenceID string, checkoutURL string, err error)
}

// IIntegrationRepository obtiene credenciales de MercadoPago desde integration_types
type IIntegrationRepository interface {
	GetMeliPagoConfig(ctx context.Context) (*MeliPagoConfig, error)
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
