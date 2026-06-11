package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

type voucherDocumentRef struct {
	ID int `json:"id"`
}

type voucherCustomer struct {
	Identification string `json:"identification"`
	BranchOffice   int    `json:"branch_office"`
}

type voucherDue struct {
	Prefix      string `json:"prefix"`
	Consecutive int    `json:"consecutive"`
	Quarter     int    `json:"quarter"`
	Year        int    `json:"year"`
}

type voucherItem struct {
	Due   voucherDue `json:"due"`
	Value float64    `json:"value"`
}

type voucherPayment struct {
	ID    int     `json:"id"`
	Value float64 `json:"value"`
}

type voucherRequest struct {
	Document     voucherDocumentRef `json:"document"`
	Date         string             `json:"date"`
	Type         string             `json:"type"`
	Customer     voucherCustomer    `json:"customer"`
	Items        []voucherItem      `json:"items"`
	Payment      voucherPayment     `json:"payment"`
	Observations string             `json:"observations,omitempty"`
}

type voucherResponse struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	Number int                  `json:"number"`
	Errors []response.SiigoError `json:"Errors,omitempty"`
}

func intFromConfig(config map[string]interface{}, keys ...string) (int, bool) {
	for _, key := range keys {
		switch v := config[key].(type) {
		case float64:
			if v > 0 {
				return int(v), true
			}
		case int:
			if v > 0 {
				return v, true
			}
		}
	}
	return 0, false
}

func (c *Client) CreateCashReceipt(ctx context.Context, req *dtos.CreateCashReceiptRequest) (*dtos.CreateCashReceiptResult, error) {
	result := &dtos.CreateCashReceiptResult{}

	c.log.Info(ctx).
		Str("invoice_number", req.InvoiceNumber).
		Msg("Creating Siigo cash receipt (voucher)")

	documentID, ok := intFromConfig(req.Config, "cash_receipt_document_id", "voucher_document_id")
	if !ok {
		return result, fmt.Errorf("cash_receipt_document_id no configurado: se requiere el id del tipo de comprobante RC en Siigo")
	}

	paymentID, ok := intFromConfig(req.Config, "cash_receipt_payment_id", "payment_method_id")
	if !ok {
		return result, fmt.Errorf("cash_receipt_payment_id no configurado: se requiere el id del medio de pago en Siigo")
	}

	token, err := c.authenticate(ctx, req.Credentials.Username, req.Credentials.AccessKey, req.Credentials.AccountID, req.Credentials.PartnerID, req.Credentials.BaseURL)
	if err != nil {
		return result, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	invoice, err := c.findInvoiceByName(ctx, token, req.Credentials, req.InvoiceNumber)
	if err != nil {
		return result, err
	}

	value := invoice.Balance
	if value <= 0 {
		return result, fmt.Errorf("la factura %s no tiene saldo pendiente en Siigo (balance %.2f)", req.InvoiceNumber, invoice.Balance)
	}

	quarter := 1
	year := time.Now().Year()
	if parsed, parseErr := time.Parse("2006-01-02", invoice.Date); parseErr == nil {
		quarter = int(parsed.Month()-1)/3 + 1
		year = parsed.Year()
	}

	voucherReq := &voucherRequest{
		Document: voucherDocumentRef{ID: documentID},
		Date:     time.Now().Format("2006-01-02"),
		Type:     "DetailedPayment",
		Customer: voucherCustomer{
			Identification: invoice.CustomerIdentification,
			BranchOffice:   invoice.CustomerBranchOffice,
		},
		Items: []voucherItem{
			{
				Due: voucherDue{
					Prefix:      invoice.Prefix,
					Consecutive: invoice.Number,
					Quarter:     quarter,
					Year:        year,
				},
				Value: value,
			},
		},
		Payment: voucherPayment{
			ID:    paymentID,
			Value: value,
		},
		Observations: "Recibo de caja generado desde Probability para " + req.InvoiceNumber,
	}

	endpoint := c.endpointURL(req.Credentials.BaseURL, "/v1/vouchers")

	var voucherResp voucherResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", req.Credentials.PartnerID).
		SetBody(voucherReq).
		SetResult(&voucherResp).
		Post(endpoint)

	result.AuditData = &dtos.AuditData{
		RequestURL:     endpoint,
		RequestPayload: voucherReq,
	}
	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo voucher request failed - network error")
		return result, fmt.Errorf("error de red al crear recibo de caja en Siigo: %w", err)
	}

	if len(voucherResp.Errors) > 0 {
		errMsg := voucherResp.Errors[0].Message
		c.log.Error(ctx).
			Str("error_code", voucherResp.Errors[0].Code).
			Str("error_msg", errMsg).
			Msg("Siigo returned business error for voucher")
		return result, fmt.Errorf("Siigo rechazo el recibo de caja: %s", errMsg)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo voucher creation failed")
		return result, fmt.Errorf("error al crear recibo de caja en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	if voucherResp.ID == "" && voucherResp.Name == "" {
		return result, fmt.Errorf("Siigo no retorno datos del recibo de caja creado")
	}

	result.ReceiptID = voucherResp.ID
	result.ReceiptName = voucherResp.Name
	result.ProviderInfo = map[string]interface{}{
		"receipt_id":     voucherResp.ID,
		"receipt_name":   voucherResp.Name,
		"receipt_number": voucherResp.Number,
		"invoice_number": req.InvoiceNumber,
		"value":          value,
	}

	c.log.Info(ctx).
		Str("receipt_id", voucherResp.ID).
		Str("receipt_name", voucherResp.Name).
		Msg("Siigo cash receipt created successfully")

	return result, nil
}

func (c *Client) findInvoiceByName(ctx context.Context, token string, credentials dtos.Credentials, invoiceName string) (*dtos.InvoiceDetail, error) {
	var listResp invoiceLookupResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetQueryParam("name", invoiceName).
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/invoices"))

	if err != nil {
		return nil, fmt.Errorf("error de red al buscar factura %s en Siigo: %w", invoiceName, err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error al buscar factura %s en Siigo (codigo %d)", invoiceName, resp.StatusCode())
	}

	for _, r := range listResp.Results {
		if r.Name == invoiceName {
			return &dtos.InvoiceDetail{
				ID:                     r.ID,
				Name:                   r.Name,
				Prefix:                 r.Prefix,
				Number:                 r.Number,
				Date:                   r.Date,
				CustomerIdentification: r.Customer.Identification,
				CustomerBranchOffice:   r.Customer.BranchOffice,
				Total:                  r.Total,
				Balance:                r.Balance,
			}, nil
		}
	}

	return nil, fmt.Errorf("factura %s no encontrada en Siigo", invoiceName)
}

type invoiceLookupResponse struct {
	Results []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Prefix   string `json:"prefix"`
		Number   int    `json:"number"`
		Date     string `json:"date"`
		Customer struct {
			Identification string `json:"identification"`
			BranchOffice   int    `json:"branch_office"`
		} `json:"customer"`
		Total   float64 `json:"total"`
		Balance float64 `json:"balance"`
	} `json:"results"`
}
