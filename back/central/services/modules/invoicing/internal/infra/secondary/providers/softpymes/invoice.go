package softpymes

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes/response"
)

// CreateInvoice crea una factura electrónica en Softpymes
func (c *Client) CreateInvoice(ctx context.Context, token string, request *ports.InvoiceRequest) (*ports.InvoiceResponse, error) {
	c.log.Info(ctx).
		Str("order_id", request.Invoice.OrderID).
		Msg("Creating invoice in Softpymes")

	// Convertir request a formato Softpymes
	softpymesReq := mappers.ToInvoiceRequest(request)

	var invoiceResp response.InvoiceResponse

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(softpymesReq).
		SetResult(&invoiceResp).
		Post("/sales_invoice/")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create invoice")
		return nil, fmt.Errorf("invoice creation request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", invoiceResp.Error).
			Msg("Invoice creation failed")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return nil, fmt.Errorf("authentication token expired")
		}

		return nil, fmt.Errorf("invoice creation failed: %s", invoiceResp.Error)
	}

	// Verificar respuesta
	if !invoiceResp.Success {
		c.log.Error(ctx).
			Str("message", invoiceResp.Message).
			Msg("Invoice creation unsuccessful")
		return nil, fmt.Errorf("invoice creation unsuccessful: %s", invoiceResp.Message)
	}

	// Convertir respuesta a formato de dominio
	result := mappers.FromInvoiceResponse(&invoiceResp)
	if result == nil {
		return nil, fmt.Errorf("failed to parse invoice response")
	}

	c.log.Info(ctx).
		Str("invoice_number", result.InvoiceNumber).
		Str("cufe", *result.CUFE).
		Msg("Invoice created successfully in Softpymes")

	return result, nil
}

// CancelInvoice cancela una factura en Softpymes
func (c *Client) CancelInvoice(ctx context.Context, token string, externalID string, reason string) error {
	c.log.Info(ctx).
		Str("external_id", externalID).
		Str("reason", reason).
		Msg("Cancelling invoice in Softpymes")

	// Request de cancelación
	cancelReq := map[string]interface{}{
		"invoice_id": externalID,
		"reason":     reason,
	}

	var cancelResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(cancelReq).
		SetResult(&cancelResp).
		Post("/sales_invoice/cancel")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to cancel invoice")
		return fmt.Errorf("invoice cancellation request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", cancelResp.Error).
			Msg("Invoice cancellation failed")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return fmt.Errorf("authentication token expired")
		}

		return fmt.Errorf("invoice cancellation failed: %s", cancelResp.Error)
	}

	// Verificar respuesta
	if !cancelResp.Success {
		return fmt.Errorf("invoice cancellation unsuccessful: %s", cancelResp.Message)
	}

	c.log.Info(ctx).
		Str("external_id", externalID).
		Msg("Invoice cancelled successfully in Softpymes")

	return nil
}

// GetInvoiceStatus consulta el estado de una factura en Softpymes
func (c *Client) GetInvoiceStatus(ctx context.Context, token string, externalID string) (string, error) {
	c.log.Info(ctx).
		Str("external_id", externalID).
		Msg("Getting invoice status from Softpymes")

	var statusResp struct {
		Success bool   `json:"success"`
		Status  string `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetQueryParam("invoice_id", externalID).
		SetResult(&statusResp).
		Get("/sales_invoice/status")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get invoice status")
		return "", fmt.Errorf("invoice status request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("error", statusResp.Error).
			Msg("Failed to get invoice status")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return "", fmt.Errorf("authentication token expired")
		}

		return "", fmt.Errorf("get invoice status failed: %s", statusResp.Error)
	}

	if !statusResp.Success {
		return "", fmt.Errorf("get invoice status unsuccessful: %s", statusResp.Message)
	}

	c.log.Info(ctx).
		Str("external_id", externalID).
		Str("status", statusResp.Status).
		Msg("Invoice status retrieved successfully")

	return statusResp.Status, nil
}
