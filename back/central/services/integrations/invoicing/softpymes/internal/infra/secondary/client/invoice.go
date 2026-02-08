package client

import (
	"context"
	"fmt"
)

// InvoiceResponse representa la respuesta de creación de factura de Softpymes
// La API retorna los datos de la factura dentro de un campo "data"
type InvoiceResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Error   string          `json:"error,omitempty"`
	Data    *InvoiceRespData `json:"data,omitempty"`
}

// InvoiceRespData contiene los datos de la factura creada por Softpymes
type InvoiceRespData struct {
	InvoiceID     string `json:"invoice_id"`
	InvoiceNumber string `json:"invoice_number"`
	CUFE          string `json:"cufe"`
	PDFURL        string `json:"pdf_url"`
	XMLURL        string `json:"xml_url"`
	InvoiceURL    string `json:"invoice_url"`
	IssuedAt      string `json:"issued_at"`
	Status        string `json:"status"`
	QRCode        string `json:"qr_code,omitempty"`
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

	// Extraer referer del config de la integración
	// El config contiene: api_url, referer, company_nit, company_name, test_mode
	config, ok := invoiceData["config"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("config not found in invoice data")
	}

	referer, ok := config["referer"].(string)
	if !ok || referer == "" {
		return fmt.Errorf("referer not found in config")
	}

	// Autenticar
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Preparar request de factura (simplificado por ahora)
	invoiceReq := map[string]interface{}{
		"customer": invoiceData["customer"],
		"items":    invoiceData["items"],
		"total":    invoiceData["total"],
		"order_id": invoiceData["order_id"],
	}

	var invoiceResp InvoiceResponse

	// Hacer llamado a la API
	requestURL := "/app/integration/sales_invoice/"
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer). // Header requerido por Softpymes
		SetBody(invoiceReq).
		SetResult(&invoiceResp).
		SetDebug(true).
		Post(requestURL) // URL correcta según documentación

	// Capturar audit data para sync logs (siempre, independiente del resultado)
	auditData := map[string]interface{}{
		"request_url":     requestURL,
		"request_payload": invoiceReq,
	}
	if resp != nil {
		auditData["response_status"] = resp.StatusCode()
		auditData["response_body"] = string(resp.Body())
	}
	invoiceData["_audit"] = auditData

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

	if invoiceResp.Data == nil {
		c.log.Warn(ctx).Msg("Invoice created but no data returned")
		return nil
	}

	c.log.Info(ctx).
		Str("invoice_number", invoiceResp.Data.InvoiceNumber).
		Str("cufe", invoiceResp.Data.CUFE).
		Str("pdf_url", invoiceResp.Data.PDFURL).
		Msg("Invoice created successfully in Softpymes")

	// Actualizar invoiceData con los datos de respuesta
	invoiceData["external_id"] = invoiceResp.Data.InvoiceID
	invoiceData["invoice_number"] = invoiceResp.Data.InvoiceNumber
	invoiceData["cufe"] = invoiceResp.Data.CUFE
	invoiceData["invoice_url"] = invoiceResp.Data.InvoiceURL
	invoiceData["pdf_url"] = invoiceResp.Data.PDFURL
	invoiceData["xml_url"] = invoiceResp.Data.XMLURL
	invoiceData["issued_at"] = invoiceResp.Data.IssuedAt

	return nil
}
