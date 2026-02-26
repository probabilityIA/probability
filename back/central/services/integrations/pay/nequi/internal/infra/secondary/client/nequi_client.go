package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	sandboxBaseURL    = "https://api.sandbox.nequi.com/payments/v2"
	productionBaseURL = "https://api.nequi.com/payments/v2"
)

// Request structures
type nequiDestination struct {
	ServiceName      string `json:"ServiceName"`
	ServiceOperation string `json:"ServiceOperation"`
	ServiceRegion    string `json:"ServiceRegion"`
	ServiceVersion   string `json:"ServiceVersion"`
}

type nequiHeader struct {
	Channel     string           `json:"Channel"`
	RequestDate string           `json:"RequestDate"`
	MessageID   string           `json:"MessageID"`
	ClientID    string           `json:"ClientID"`
	Destination nequiDestination `json:"Destination"`
}

type nequiQRRequest struct {
	Code       string `json:"code"`
	Value      string `json:"value"`
	Reference1 string `json:"reference1"`
	Reference2 string `json:"reference2"`
	Reference3 string `json:"reference3"`
}

type nequiBody struct {
	Any struct {
		GenerateCodeQRRQ nequiQRRequest `json:"generateCodeQRRQ"`
	} `json:"any"`
}

type nequiRequest struct {
	RequestMessage struct {
		RequestHeader nequiHeader `json:"RequestHeader"`
		RequestBody   nequiBody   `json:"RequestBody"`
	} `json:"RequestMessage"`
}

// Response structures
type nequiQRResponse struct {
	QrValue       string `json:"qrValue"`
	TransactionId string `json:"transactionId"`
}

type nequiResponseBody struct {
	Any struct {
		GenerateCodeQRRS nequiQRResponse `json:"generateCodeQRRS"`
	} `json:"any"`
}

type nequiResponse struct {
	ResponseMessage struct {
		ResponseHeader interface{}       `json:"ResponseHeader"`
		ResponseBody   nequiResponseBody `json:"ResponseBody"`
	} `json:"ResponseMessage"`
}

// NequiClient implementa ports.INequiClient
type NequiClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente Nequi
func New(logger log.ILogger) ports.INequiClient {
	return &NequiClient{
		log: logger.WithModule("nequi.client"),
	}
}

// GenerateQR genera un código QR de pago en Nequi
func (c *NequiClient) GenerateQR(ctx context.Context, config *ports.NequiConfig, amount float64, reference string) (string, string, error) {
	baseURL := sandboxBaseURL
	if config.Environment == "production" {
		baseURL = productionBaseURL
	}

	client := resty.New().SetBaseURL(baseURL)

	phoneCode := config.PhoneCode
	if phoneCode == "" {
		phoneCode = "NIT_1"
	}

	req := nequiRequest{}
	req.RequestMessage.RequestHeader = nequiHeader{
		Channel:     "PQR03-C001",
		RequestDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
		MessageID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		ClientID:    "12345",
		Destination: nequiDestination{
			ServiceName:      "PaymentsService",
			ServiceOperation: "generateCodeQR",
			ServiceRegion:    "C001",
			ServiceVersion:   "1.2.0",
		},
	}
	req.RequestMessage.RequestBody.Any.GenerateCodeQRRQ = nequiQRRequest{
		Code:       phoneCode,
		Value:      fmt.Sprintf("%.0f", amount),
		Reference1: reference,
		Reference2: "pay",
		Reference3: "probability",
	}

	var result nequiResponse
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("x-api-key", config.APIKey).
		SetBody(req).
		SetResult(&result).
		Post("/-services-paymentservice-generatecodeqr")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Error calling Nequi API")
		return "", "", fmt.Errorf("nequi api call failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).Str("status", resp.Status()).Str("body", resp.String()).Msg("Nequi API returned error")
		return "", "", fmt.Errorf("nequi api error: status=%s", resp.Status())
	}

	qr := result.ResponseMessage.ResponseBody.Any.GenerateCodeQRRS.QrValue
	txID := result.ResponseMessage.ResponseBody.Any.GenerateCodeQRRS.TransactionId

	if qr == "" {
		return "", "", fmt.Errorf("empty qr value in nequi response")
	}

	return qr, txID, nil
}
