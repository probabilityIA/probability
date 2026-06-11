package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

type creditNoteItem struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
}

type creditNoteRequest struct {
	Document struct {
		ID int `json:"id"`
	} `json:"document"`
	Date     string `json:"date"`
	Invoice  string `json:"invoice"`
	Customer struct {
		Identification string `json:"identification"`
	} `json:"customer"`
	Items   []creditNoteItem `json:"items"`
	Reason  string           `json:"reason,omitempty"`
	Observations string      `json:"observations,omitempty"`
}

type creditNoteResponse struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	Number   int                   `json:"number"`
	Metadata struct {
		CUFE string `json:"cufe"`
	} `json:"metadata"`
	Stamp struct {
		CUFE string `json:"cufe"`
	} `json:"stamp"`
	Errors []response.SiigoError `json:"Errors,omitempty"`
}

func (c *Client) CreateCreditNote(ctx context.Context, req *dtos.CreateCreditNoteRequest) (*dtos.CreateCreditNoteResult, error) {
	result := &dtos.CreateCreditNoteResult{}

	c.log.Info(ctx).
		Str("invoice_external_id", req.InvoiceExternalID).
		Float64("amount", req.Amount).
		Msg("Creating Siigo credit note")

	documentID, ok := intFromConfig(req.Config, "credit_note_document_id")
	if !ok {
		return result, fmt.Errorf("credit_note_document_id no configurado: se requiere el id del tipo de documento NC en Siigo")
	}

	token, err := c.authenticate(ctx, req.Credentials.Username, req.Credentials.AccessKey, req.Credentials.AccountID, req.Credentials.PartnerID, req.Credentials.BaseURL)
	if err != nil {
		return result, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	body := &creditNoteRequest{}
	body.Document.ID = documentID
	body.Date = time.Now().Format("2006-01-02")
	body.Invoice = req.InvoiceExternalID
	body.Customer.Identification = req.CustomerDNI
	body.Items = []creditNoteItem{
		{
			Code:        "NC",
			Description: reasonOrDefault(req.Reason),
			Quantity:    1,
			Price:       req.Amount,
		},
	}
	body.Reason = req.Reason
	body.Observations = "Nota de credito generada desde Probability para " + req.InvoiceNumber

	endpoint := c.endpointURL(req.Credentials.BaseURL, "/v1/credit-notes")

	var ncResp creditNoteResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", req.Credentials.PartnerID).
		SetBody(body).
		SetResult(&ncResp).
		Post(endpoint)

	result.AuditData = &dtos.AuditData{
		RequestURL:     endpoint,
		RequestPayload: body,
	}
	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo credit note request failed - network error")
		return result, fmt.Errorf("error de red al crear nota de credito en Siigo: %w", err)
	}

	if len(ncResp.Errors) > 0 {
		errMsg := ncResp.Errors[0].Message
		c.log.Error(ctx).Str("error_msg", errMsg).Msg("Siigo returned business error for credit note")
		return result, fmt.Errorf("Siigo rechazo la nota de credito: %s", errMsg)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo credit note creation failed")
		return result, fmt.Errorf("error al crear nota de credito en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	if ncResp.ID == "" && ncResp.Name == "" {
		return result, fmt.Errorf("Siigo no retorno datos de la nota de credito creada")
	}

	cufe := ncResp.Stamp.CUFE
	if cufe == "" {
		cufe = ncResp.Metadata.CUFE
	}

	result.CreditNoteID = ncResp.ID
	result.CreditNoteNumber = ncResp.Name
	result.CUFE = cufe
	result.ProviderInfo = map[string]interface{}{
		"credit_note_id":     ncResp.ID,
		"credit_note_name":   ncResp.Name,
		"credit_note_number": ncResp.Number,
		"invoice_number":     req.InvoiceNumber,
	}

	c.log.Info(ctx).
		Str("credit_note_id", ncResp.ID).
		Str("credit_note_name", ncResp.Name).
		Msg("Siigo credit note created successfully")

	return result, nil
}

func reasonOrDefault(reason string) string {
	if reason == "" {
		return "Anulacion / devolucion"
	}
	return reason
}
