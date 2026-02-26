package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const boldBaseURL = "https://integrations.api.bold.co"

// Estructuras de request Bold
type boldAmount struct {
	Currency    string        `json:"currency"`
	TotalAmount float64       `json:"total_amount"`
	Taxes       []interface{} `json:"taxes"`
}

type boldCreateLinkRequest struct {
	AmountType     string     `json:"amount_type"`
	Amount         boldAmount `json:"amount"`
	Description    string     `json:"description"`
	PaymentMethods []string   `json:"payment_methods"`
}

// Estructuras de response Bold
type boldPayload struct {
	PaymentLink string `json:"payment_link"`
	URL         string `json:"url"`
}

type boldCreateLinkResponse struct {
	Payload boldPayload   `json:"payload"`
	Errors  []interface{} `json:"errors"`
}

// BoldClient implementa ports.IBoldClient
type BoldClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente Bold
func New(logger log.ILogger) ports.IBoldClient {
	return &BoldClient{
		log: logger.WithModule("bold.client"),
	}
}

// CreatePaymentLink crea un link de pago en Bold
func (c *BoldClient) CreatePaymentLink(ctx context.Context, config *ports.BoldConfig, amount float64, currency, reference, description string) (string, string, error) {
	if currency == "" {
		currency = "COP"
	}

	reqBody := boldCreateLinkRequest{
		AmountType: "CLOSE",
		Amount: boldAmount{
			Currency:    currency,
			TotalAmount: amount,
			Taxes:       []interface{}{},
		},
		Description: description,
		PaymentMethods: []string{
			"CREDIT_CARD",
			"PSE",
			"NEQUI",
			"BOTON_BANCOLOMBIA",
		},
	}

	var result boldCreateLinkResponse
	client := resty.New().SetBaseURL(boldBaseURL)

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("x-api-key %s", config.APIKey)).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&result).
		Post("/online/link/v1")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Error calling Bold API")
		return "", "", fmt.Errorf("bold api call failed: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).Str("status", resp.Status()).Str("body", resp.String()).Msg("Bold API returned error")
		return "", "", fmt.Errorf("bold api error: status=%s body=%s", resp.Status(), resp.String())
	}

	if len(result.Errors) > 0 {
		c.log.Error(ctx).Interface("errors", result.Errors).Msg("Bold API returned errors in payload")
		return "", "", fmt.Errorf("bold api returned errors: %v", result.Errors)
	}

	linkID := result.Payload.PaymentLink
	checkoutURL := result.Payload.URL

	if linkID == "" {
		return "", "", fmt.Errorf("empty payment_link in bold response")
	}

	return linkID, checkoutURL, nil
}
