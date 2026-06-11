package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

func (c *Client) AnnulInvoice(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.AnnulInvoiceResult, error) {
	result := &dtos.AnnulInvoiceResult{}

	c.log.Info(ctx).
		Str("siigo_invoice_id", invoiceID).
		Msg("Annulling Siigo invoice")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return result, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	endpoint := c.endpointURL(credentials.BaseURL, "/v1/invoices/"+invoiceID+"/annul")

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		Post(endpoint)

	result.AuditData = &dtos.AuditData{
		RequestURL: endpoint,
	}
	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo annul request failed - network error")
		return result, fmt.Errorf("error de red al anular factura en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo annul invoice failed")
		switch resp.StatusCode() {
		case 404:
			return result, fmt.Errorf("factura %s no encontrada en Siigo", invoiceID)
		case 409:
			return result, fmt.Errorf("Siigo no permite anular esta factura (annul_not_allowed): puede tener recibos de caja u otros documentos asociados")
		default:
			return result, fmt.Errorf("error al anular factura en Siigo (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
		}
	}

	c.log.Info(ctx).
		Str("siigo_invoice_id", invoiceID).
		Msg("Siigo invoice annulled successfully")

	return result, nil
}
