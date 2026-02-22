package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/response"
)

// CreateInvoice crea una factura electrónica en Factus
// Endpoint: POST /v1/bills/validate
func (c *Client) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
	result := &dtos.CreateInvoiceResult{}

	if req.Credentials.ClientID == "" {
		return result, fmt.Errorf("factus create_invoice: client_id missing")
	}

	token, err := c.authenticate(
		ctx,
		req.Credentials.BaseURL,
		req.Credentials.ClientID,
		req.Credentials.ClientSecret,
		req.Credentials.Username,
		req.Credentials.Password,
	)
	if err != nil {
		return result, fmt.Errorf("factus create_invoice: authentication failed: %w", err)
	}

	billReq := mappers.BuildCreateBillRequest(req)

	c.log.Info(ctx).
		Str("order_id", req.OrderID).
		Str("reference_code", billReq.ReferenceCode).
		Int("numbering_range_id", billReq.NumberingRangeID).
		Int("items_count", len(billReq.Items)).
		Msg("📄 Creating Factus invoice")

	result.AuditData = &dtos.AuditData{
		RequestURL:     "/v1/bills/validate",
		RequestPayload: billReq,
	}

	var apiResp response.CreateBill

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(billReq).
		SetResult(&apiResp).
		Post(c.endpointURL(req.Credentials.BaseURL, "/v1/bills/validate"))
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Factus create_invoice request failed - network error")
		return result, fmt.Errorf("factus create_invoice request failed: %w", err)
	}

	result.AuditData.ResponseStatus = resp.StatusCode()
	result.AuditData.ResponseBody = string(resp.Body())

	if resp.IsError() {
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return result, fmt.Errorf("factus create_invoice: authentication token expired (401)")
		}
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Factus create_invoice failed")
		return result, fmt.Errorf("factus create_invoice failed (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	bill := apiResp.Data.Bill
	result.InvoiceNumber = bill.Number
	result.ExternalID = strconv.Itoa(bill.ID)
	result.CUFE = bill.CUFE
	result.QRCode = bill.QR
	result.Total = bill.Total
	result.IssuedAt = bill.Validated

	c.log.Info(ctx).
		Str("invoice_number", result.InvoiceNumber).
		Str("cufe", result.CUFE).
		Str("total", result.Total).
		Msg("✅ Factus invoice created successfully")

	return result, nil
}
