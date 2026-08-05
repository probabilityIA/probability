package ports

import (
	"context"
)

type BancolombiaConfig struct {
	ClientID       string
	ClientSecret   string
	WebhookSecret  string
	Environment    string
	BaseURL        string
	TokenPath      string
	QRGeneratePath string
	QRQueryPath    string
	Scope          string
	MerchantID     string
	MerchantName   string
}

func (c *BancolombiaConfig) WebhookSecretKey() (string, bool) {
	if c == nil || c.WebhookSecret == "" {
		return "", false
	}
	return c.WebhookSecret, true
}

type QRGenerationResult struct {
	ExternalID    string
	QRValue       string
	QRImageBase64 string
	RawResponse   map[string]interface{}
}

type IBancolombiaClient interface {
	GenerateQR(ctx context.Context, config *BancolombiaConfig, amount float64, currency, reference, description string) (*QRGenerationResult, error)
}

type IIntegrationRepository interface {
	GetConfig(ctx context.Context) (*BancolombiaConfig, error)
	GetConfigForMode(ctx context.Context, testMode bool) (*BancolombiaConfig, error)
}

type IResponsePublisher interface {
	PublishPaymentResponse(ctx context.Context, msg *PaymentResponseMsg) error
}

type PaymentResponseMsg struct {
	PaymentTransactionID uint                   `json:"payment_transaction_id"`
	GatewayCode          string                 `json:"gateway_code"`
	Status               string                 `json:"status"`
	ExternalID           *string                `json:"external_id,omitempty"`
	GatewayResponse      map[string]interface{} `json:"gateway_response,omitempty"`
	Error                string                 `json:"error,omitempty"`
	ErrorCode            string                 `json:"error_code,omitempty"`
	CorrelationID        string                 `json:"correlation_id"`
	ProcessingTimeMs     int64                  `json:"processing_time_ms"`
}
