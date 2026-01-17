package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/modules/wallet/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type NequiService struct {
	client *resty.Client
	config env.IConfig
	logger log.ILogger
}

func NewNequiService(config env.IConfig, logger log.ILogger) domain.INequiService {
	client := resty.New()
	client.SetBaseURL("https://api.sandbox.nequi.com/payments/v2")
	// client.SetDebug(true)
	return &NequiService{
		client: client,
		config: config,
		logger: logger,
	}
}

// Request Structures
type RequestHeader struct {
	Channel     string      `json:"Channel"`
	RequestDate string      `json:"RequestDate"`
	MessageID   string      `json:"MessageID"`
	ClientID    string      `json:"ClientID"`
	Destination Destination `json:"Destination"`
}

type Destination struct {
	ServiceName      string `json:"ServiceName"`
	ServiceOperation string `json:"ServiceOperation"`
	ServiceRegion    string `json:"ServiceRegion"`
	ServiceVersion   string `json:"ServiceVersion"`
}

type GenerateCodeQRRQ struct {
	Code       string `json:"code"`
	Value      string `json:"value"`
	Reference1 string `json:"reference1"`
	Reference2 string `json:"reference2"`
	Reference3 string `json:"reference3"`
}

type RequestBody struct {
	Any struct {
		GenerateCodeQRRQ GenerateCodeQRRQ `json:"generateCodeQRRQ"`
	} `json:"any"`
}

type RequestMessage struct {
	RequestHeader RequestHeader `json:"RequestHeader"`
	RequestBody   RequestBody   `json:"RequestBody"`
}

type NequiRequest struct {
	RequestMessage RequestMessage `json:"RequestMessage"`
}

// Response Structures
type GenerateCodeQRRS struct {
	QrValue       string `json:"qrValue"`
	TransactionId string `json:"transactionId"`
}

type ResponseBody struct {
	Any struct {
		GenerateCodeQRRS GenerateCodeQRRS `json:"generateCodeQRRS"`
	} `json:"any"`
}

type ResponseMessage struct {
	ResponseHeader interface{}  `json:"ResponseHeader"`
	ResponseBody   ResponseBody `json:"ResponseBody"`
}

type NequiResponse struct {
	ResponseMessage ResponseMessage `json:"ResponseMessage"`
}

func (s *NequiService) GenerateQR(ctx context.Context, amount float64) (string, string, error) {
	apiKey := s.config.Get("NEQUI_API_KEY")
	// token := s.config.Get("NEQUI_ACCESS_TOKEN")

	reqBody := NequiRequest{
		RequestMessage: RequestMessage{
			RequestHeader: RequestHeader{
				Channel:     "PQR03-C001",
				RequestDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
				MessageID:   fmt.Sprintf("%d", time.Now().UnixNano()),
				ClientID:    "12345",
				Destination: Destination{
					ServiceName:      "PaymentsService",
					ServiceOperation: "generateCodeQR",
					ServiceRegion:    "C001",
					ServiceVersion:   "1.2.0",
				},
			},
			RequestBody: RequestBody{
				Any: struct {
					GenerateCodeQRRQ GenerateCodeQRRQ `json:"generateCodeQRRQ"`
				}{
					GenerateCodeQRRQ: GenerateCodeQRRQ{
						Code:       "NIT_1", // Should be configured?
						Value:      fmt.Sprintf("%.0f", amount),
						Reference1: "Ref1",
						Reference2: "Ref2",
						Reference3: "Ref3",
					},
				},
			},
		},
	}

	var result NequiResponse
	// Path: /-services-paymentservice-generatecodeqr
	resp, err := s.client.R().
		SetContext(ctx).
		SetHeader("x-api-key", apiKey).
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(reqBody).
		SetResult(&result).
		Post("/-services-paymentservice-generatecodeqr")

	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Error calling Nequi API")
		return "", "", err
	}

	if resp.IsError() {
		s.logger.Error(ctx).Err(fmt.Errorf("status: %s, body: %s", resp.Status(), resp.String())).Msg("Nequi API returned error")
		return "", "", fmt.Errorf("nequi api error: %s", resp.String())
	}

	qr := result.ResponseMessage.ResponseBody.Any.GenerateCodeQRRS.QrValue
	txId := result.ResponseMessage.ResponseBody.Any.GenerateCodeQRRS.TransactionId

	if qr == "" {
		return "", "", fmt.Errorf("empty qr code in response")
	}

	return qr, txId, nil
}
