package client

import (
	"context"
	"fmt"
)

// InvoiceResponse representa la respuesta de creación de factura de Softpymes
<<<<<<< HEAD
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
=======
// Según documentación oficial: https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostSaleInvoice
type InvoiceResponse struct {
	Message string      `json:"message"` // "Se ha creado la factura de venta en Pymes+ correctamente!"
	Info    *InvoiceInfo `json:"info,omitempty"`
}

// InvoiceInfo contiene los datos de la factura creada por Softpymes
type InvoiceInfo struct {
	Date           string  `json:"date"`           // "2023-10-25T10:39:13.000Z"
	DocumentNumber string  `json:"documentNumber"` // "ABC0000000000"
	Subtotal       float64 `json:"subtotal"`
	Discount       float64 `json:"discount"`
	IVA            float64 `json:"iva"`
	Withholding    float64 `json:"withholding"`
	Total          float64 `json:"total"`
	DocsFe         *DocsFe `json:"docsFe,omitempty"`
}

// DocsFe contiene información de validación de la factura electrónica
type DocsFe struct {
	Status  bool   `json:"status"`  // true = válido
	Message string `json:"message"` // "Documento válido enviado al proveedor tecnológico"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
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

<<<<<<< HEAD
	// Extraer referer del config
=======
	// Extraer referer del config de la integración
	// El config contiene: api_url, referer, company_nit, company_name, test_mode
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
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
<<<<<<< HEAD
		"customer":    invoiceData["customer"],
		"items":       invoiceData["items"],
		"total":       invoiceData["total"],
		"order_id":    invoiceData["order_id"],
=======
		"customer": invoiceData["customer"],
		"items":    invoiceData["items"],
		"total":    invoiceData["total"],
		"order_id": invoiceData["order_id"],
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	}

	var invoiceResp InvoiceResponse

	// Hacer llamado a la API
<<<<<<< HEAD
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(invoiceReq).
		SetResult(&invoiceResp).
		Post("/sales_invoice/")
=======
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create invoice")
		return fmt.Errorf("invoice creation request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
<<<<<<< HEAD
			Str("error", invoiceResp.Error).
=======
			Str("response", string(resp.Body())).
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
			Msg("Invoice creation failed")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return fmt.Errorf("authentication token expired")
		}

<<<<<<< HEAD
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
=======
		return fmt.Errorf("invoice creation failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Verificar que haya info en la respuesta
	if invoiceResp.Info == nil {
		c.log.Warn(ctx).
			Str("message", invoiceResp.Message).
			Msg("Invoice response has no info")
		return fmt.Errorf("invoice response has no info: %s", invoiceResp.Message)
	}

	c.log.Info(ctx).
		Str("document_number", invoiceResp.Info.DocumentNumber).
		Str("date", invoiceResp.Info.Date).
		Float64("total", invoiceResp.Info.Total).
		Str("message", invoiceResp.Message).
		Msg("Invoice created successfully in Softpymes")

	// Actualizar invoiceData con los datos de respuesta
	// Nota: Softpymes retorna el documentNumber pero no un ID único de factura
	invoiceData["external_id"] = invoiceResp.Info.DocumentNumber // Usar documentNumber como ID
	invoiceData["invoice_number"] = invoiceResp.Info.DocumentNumber
	invoiceData["issued_at"] = invoiceResp.Info.Date

	// Información adicional del provider
	providerInfo := map[string]interface{}{
		"subtotal":    invoiceResp.Info.Subtotal,
		"discount":    invoiceResp.Info.Discount,
		"iva":         invoiceResp.Info.IVA,
		"withholding": invoiceResp.Info.Withholding,
		"total":       invoiceResp.Info.Total,
	}

	if invoiceResp.Info.DocsFe != nil {
		providerInfo["dian_status"] = invoiceResp.Info.DocsFe.Status
		providerInfo["dian_message"] = invoiceResp.Info.DocsFe.Message
	}

	invoiceData["provider_info"] = providerInfo
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	return nil
}
