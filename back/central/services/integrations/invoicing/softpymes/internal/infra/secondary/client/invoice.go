package client

import (
	"context"
	"fmt"
)

// InvoiceResponse representa la respuesta de creación de factura de Softpymes
type InvoiceResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	Error         string `json:"error"`
	InvoiceNumber string `json:"invoice_number"`
	ExternalID    string `json:"external_id"`
	InvoiceURL    string `json:"invoice_url"`
	PDFURL        string `json:"pdf_url"`
	XMLURL        string `json:"xml_url"`
	CUFE          string `json:"cufe"`
	IssuedAt      string `json:"issued_at"`
}

// CreateInvoice crea una factura electrónica en Softpymes
func (c *Client) CreateInvoice(ctx context.Context, invoiceData map[string]interface{}) error {
	c.log.Info(ctx).Interface("data", invoiceData).Msg("Creating invoice in Softpymes")

	// Extraer credenciales del map
	credentials, ok := invoiceData["credentials"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("credentials not found in invoice data")
	}

	apiKey, ok := credentials["api_key"].(string)
	if !ok || apiKey == "" {
		return fmt.Errorf("api_key not found in credentials")
	}

	apiSecret, ok := credentials["api_secret"].(string)
	if !ok || apiSecret == "" {
		return fmt.Errorf("api_secret not found in credentials")
	}

	// Autenticar
	token, err := c.authenticate(ctx, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Preparar request de factura (simplificado por ahora)
	invoiceReq := map[string]interface{}{
		"customer":    invoiceData["customer"],
		"items":       invoiceData["items"],
		"total":       invoiceData["total"],
		"order_id":    invoiceData["order_id"],
	}

	var invoiceResp InvoiceResponse

	// Hacer llamado a la API
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(invoiceReq).
		SetResult(&invoiceResp).
		Post("/sales_invoice/")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create invoice")
		return fmt.Errorf("invoice creation request failed: %w", err)
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
			return fmt.Errorf("authentication token expired")
		}

		return fmt.Errorf("invoice creation failed: %s", invoiceResp.Error)
	}

	// Verificar respuesta
	if !invoiceResp.Success {
		c.log.Error(ctx).
			Str("message", invoiceResp.Message).
			Msg("Invoice creation unsuccessful")
		return fmt.Errorf("invoice creation unsuccessful: %s", invoiceResp.Message)
	}

	c.log.Info(ctx).
		Str("invoice_number", invoiceResp.InvoiceNumber).
		Str("cufe", invoiceResp.CUFE).
		Msg("Invoice created successfully in Softpymes")

	// Actualizar invoiceData con los datos de respuesta
	invoiceData["external_id"] = invoiceResp.ExternalID
	invoiceData["invoice_number"] = invoiceResp.InvoiceNumber
	invoiceData["cufe"] = invoiceResp.CUFE
	invoiceData["invoice_url"] = invoiceResp.InvoiceURL
	invoiceData["pdf_url"] = invoiceResp.PDFURL
	invoiceData["xml_url"] = invoiceResp.XMLURL

	return nil
}
