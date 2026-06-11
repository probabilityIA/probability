package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type invoiceDetailResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Number   int    `json:"number"`
	Date     string `json:"date"`
	Status   string `json:"status"`
	Customer struct {
		ID             string `json:"id"`
		Identification string `json:"identification"`
		BranchOffice   int    `json:"branch_office"`
	} `json:"customer"`
	Total     float64 `json:"total"`
	Balance   float64 `json:"balance"`
	PublicURL string  `json:"public_url"`
	Stamp     struct {
		Status string `json:"status"`
		CUFE   string `json:"cufe"`
	} `json:"stamp"`
	Metadata struct {
		CUFE string `json:"cufe"`
	} `json:"metadata"`
}

func (c *Client) GetInvoiceByID(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.InvoiceDetail, error) {
	c.log.Info(ctx).
		Str("siigo_invoice_id", invoiceID).
		Msg("Getting Siigo invoice by id")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	var detail invoiceDetailResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetResult(&detail).
		Get(c.endpointURL(credentials.BaseURL, "/v1/invoices/"+invoiceID))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo get invoice request failed - network error")
		return nil, fmt.Errorf("error de red al consultar factura en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo get invoice failed")
		if resp.StatusCode() == 404 {
			return nil, fmt.Errorf("factura %s no encontrada en Siigo", invoiceID)
		}
		return nil, fmt.Errorf("error al consultar factura en Siigo (codigo %d)", resp.StatusCode())
	}

	cufe := detail.Stamp.CUFE
	if cufe == "" {
		cufe = detail.Metadata.CUFE
	}

	return &dtos.InvoiceDetail{
		ID:                     detail.ID,
		Name:                   detail.Name,
		Prefix:                 detail.Prefix,
		Number:                 detail.Number,
		Date:                   detail.Date,
		CustomerID:             detail.Customer.ID,
		CustomerIdentification: detail.Customer.Identification,
		CustomerBranchOffice:   detail.Customer.BranchOffice,
		Total:                  detail.Total,
		Balance:                detail.Balance,
		Status:                 detail.Status,
		StampStatus:            detail.Stamp.Status,
		CUFE:                   cufe,
		PublicURL:              detail.PublicURL,
	}, nil
}

func (c *Client) GetStampErrors(ctx context.Context, credentials dtos.Credentials, invoiceID string) ([]dtos.StampError, error) {
	c.log.Info(ctx).
		Str("siigo_invoice_id", invoiceID).
		Msg("Getting Siigo stamp errors")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		Get(c.endpointURL(credentials.BaseURL, "/v1/invoices/"+invoiceID+"/stamp/errors"))

	if err != nil {
		return nil, fmt.Errorf("error de red al consultar errores de timbrado en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Warn(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo stamp errors request failed")
		return nil, fmt.Errorf("error al consultar errores de timbrado en Siigo (codigo %d)", resp.StatusCode())
	}

	return parseStampErrors(resp.Body()), nil
}

type stampErrorItem struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func parseStampErrors(body []byte) []dtos.StampError {
	var wrapped struct {
		Errors []stampErrorItem `json:"Errors"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Errors) > 0 {
		return mapStampErrors(wrapped.Errors)
	}

	var plain []stampErrorItem
	if err := json.Unmarshal(body, &plain); err == nil {
		return mapStampErrors(plain)
	}

	return nil
}

func mapStampErrors(items []stampErrorItem) []dtos.StampError {
	result := make([]dtos.StampError, 0, len(items))
	for _, e := range items {
		result = append(result, dtos.StampError{Code: e.Code, Message: e.Message})
	}
	return result
}
