package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
)

// InvoiceResponse representa la respuesta de creación de factura de Softpymes
// Según documentación oficial: https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostSaleInvoice
type InvoiceResponse struct {
	Message string       `json:"message"` // "Se ha creado la factura de venta en Pymes+ correctamente!"
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
// NOTA: Softpymes retorna status como string ("Aceptado", "Pendiente", "Rechazado"), NO bool.
type DocsFe struct {
	Status  string `json:"status"`          // "Aceptado", "Pendiente", "Rechazado"
	Message string `json:"message"`         // "Documento válido enviado al proveedor tecnológico"
	Error   string `json:"error,omitempty"` // Error de la DIAN (ej: "Empresa no habilitada...")
}

// CreateInvoice crea una factura electrónica en Softpymes
// baseURL: URL base efectiva (producción o testing); vacío usa la URL del constructor
func (c *Client) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest, baseURL string) (*dtos.CreateInvoiceResult, error) {
	result := &dtos.CreateInvoiceResult{}

	// Validar credenciales
	if req.Credentials.APIKey == "" {
		return result, fmt.Errorf("api_key not found in credentials")
	}
	if req.Credentials.APISecret == "" {
		return result, fmt.Errorf("api_secret not found in credentials")
	}

	// Extraer referer del config
	referer, _ := req.Config["referer"].(string)
	if referer == "" {
		return result, fmt.Errorf("referer not found in config")
	}

	// Autenticar usando la URL efectiva
	token, err := c.authenticate(ctx, req.Credentials.APIKey, req.Credentials.APISecret, referer, baseURL)
	if err != nil {
		return result, fmt.Errorf("authentication failed: %w", err)
	}

	// Verificar idempotencia: consultar si ya existe una factura para esta orden
	if req.OrderID != "" {
		existing, err := c.findExistingInvoiceByOrderID(ctx, req.Credentials.APIKey, req.Credentials.APISecret, referer, req.OrderID, baseURL)
		if err != nil {
			// No bloqueamos la creación si la consulta falla — solo advertimos
			c.log.Warn(ctx).Err(err).
				Str("order_id", req.OrderID).
				Msg("Could not check for existing invoice, proceeding with creation")
		} else if existing != nil {
			c.log.Info(ctx).
				Str("order_id", req.OrderID).
				Str("document_number", existing.DocumentNumber).
				Msg("Invoice already exists for this order in Softpymes, skipping creation")
			result.InvoiceNumber = existing.DocumentNumber
			result.ExternalID = existing.DocumentNumber
			result.IssuedAt = existing.DocumentDate
			result.ProviderInfo = map[string]interface{}{
				"already_existed": true,
				"document_number": existing.DocumentNumber,
			}
			return result, nil
		}
	}

	// Determinar customerNit
	customerNit := req.Customer.DNI
	if customerNit == "" {
		if defaultNit, ok := req.Config["default_customer_nit"].(string); ok && defaultNit != "" {
			customerNit = defaultNit
			c.log.Info(ctx).
				Str("default_nit", defaultNit).
				Msg("Using default customer NIT from config")
		} else {
			c.log.Error(ctx).
				Msg("customerNit is required but not provided. Configure default_customer_nit in integration config or ensure customers have DNI")
			return result, fmt.Errorf("customerNit is required: customer has no DNI and no default_customer_nit configured")
		}
	}

	// Asegurar que el cliente existe en Softpymes antes de facturar
	// Retorna el branchCode real asignado por Softpymes
	customerBranch := ""
	if branch, err := c.ensureCustomerExists(ctx, token, referer, customerNit, &req.Customer, req.Config, baseURL); err != nil {
		c.log.Warn(ctx).Err(err).
			Str("customer_nit", customerNit).
			Msg("Could not ensure customer exists in Softpymes, proceeding anyway")
	} else {
		customerBranch = branch
	}

	// Fallback: usar config solo si no se pudo obtener branchCode del cliente
	if customerBranch == "" {
		customerBranch = "001" // default
		if cb, ok := req.Config["customer_branch_code"].(string); ok && cb != "" {
			customerBranch = cb
		}
	}

	// Obtener branch_code del config (default "001")
	branchCode := "001"
	if branch, ok := req.Config["branch_code"].(string); ok && branch != "" {
		branchCode = branch
	}

	// Obtener seller_nit del config (OPCIONAL)
	sellerNit := ""
	if seller, ok := req.Config["seller_nit"].(string); ok && seller != "" {
		sellerNit = seller
	}

	// Obtener resolution_id del config
	resolutionID := 0
	if resID, ok := req.Config["resolution_id"].(float64); ok {
		resolutionID = int(resID)
	} else if resID, ok := req.Config["resolution_id"].(int); ok {
		resolutionID = resID
	}

	if resolutionID == 0 {
		c.log.Warn(ctx).
			Msg("resolutionId is 0 - Softpymes may reject this invoice. Configure a valid resolution_id in integration config")
	}

	// Validar items
	if len(req.Items) == 0 {
		c.log.Error(ctx).Msg("No items found in invoice request - cannot create invoice without items")
		return result, fmt.Errorf("no items found in invoice data")
	}

	// Mapear items al formato de Softpymes
	softpymesItems := make([]map[string]interface{}, 0, len(req.Items))
	for _, item := range req.Items {
		itemCode := item.SKU
		if itemCode == "" && item.ProductID != nil {
			itemCode = *item.ProductID
		}

		quantity := float64(item.Quantity)
		discount := item.Discount

		// unitCode por defecto "UNI" (UNIDADES) - código estándar en Softpymes
		unitCode := "UNI"

		softpymesItem := map[string]interface{}{
			"itemCode": itemCode,
			"quantity": quantity,
			"discount": discount,
			"unitCode": unitCode,
		}

		// Solo incluir unitValue si tiene valor (Softpymes espera String)
		if item.UnitPrice > 0 {
			softpymesItem["unitValue"] = fmt.Sprintf("%.2f", item.UnitPrice)
		}

		softpymesItems = append(softpymesItems, softpymesItem)
	}

	// Mapear currency al formato de Softpymes
	rawCurrency := req.Currency
	if rawCurrency == "" {
		rawCurrency = "COP"
	}
	currency := mapCurrencyToSoftpymes(rawCurrency)

	// Generar documentDate en formato YYYY-MM-DD (zona horaria Colombia UTC-5)
	loc, _ := time.LoadLocation("America/Bogota")
	documentDate := time.Now().In(loc).Format("2006-01-02")

	// comment: identifica la orden de origen para idempotencia y trazabilidad
	comment := ""
	if req.OrderID != "" {
		comment = "order:" + req.OrderID
	}

	// Construir request según formato de Softpymes
	invoiceReq := map[string]interface{}{
		"documentDate":   documentDate,
		"currencyCode":   currency,
		"exchangeRate":   1.0,
		"branchCode":     branchCode,
		"customerBranch": customerBranch,
		"customerNit":    customerNit,
		"termDays":       0,
		"resolutionId":   resolutionID,
		"comment":        comment,
		"items":          softpymesItems,
	}

	// Solo incluir sellerNit si está configurado
	if sellerNit != "" {
		invoiceReq["sellerNit"] = sellerNit
	}

	c.log.Info(ctx).
		Str("document_date", documentDate).
		Str("currency", currency).
		Str("customer_nit", customerNit).
		Str("customer_branch", customerBranch).
		Str("seller_nit", sellerNit).
		Str("branch_code", branchCode).
		Int("resolution_id", resolutionID).
		Int("items_count", len(softpymesItems)).
		Interface("items", softpymesItems).
		Msg("Sending invoice request to Softpymes")

	var invoiceResp InvoiceResponse

	requestURL := c.resolveURL(baseURL, "/app/integration/sales_invoice/")
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(invoiceReq).
		SetResult(&invoiceResp).
		SetDebug(true).
		Post(requestURL)

	// Capturar audit data (siempre, independiente del resultado)
	result.AuditData = &dtos.AuditData{
		RequestURL:     requestURL,
		RequestPayload: invoiceReq,
	}
	if resp != nil {
		result.AuditData.ResponseStatus = resp.StatusCode()
		result.AuditData.ResponseBody = string(resp.Body())
	}

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to create invoice")
		return result, fmt.Errorf("invoice creation request failed: %w", err)
	}

	// Manejar errores HTTP
	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("Invoice creation failed")

		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return result, fmt.Errorf("authentication token expired")
		}

		return result, fmt.Errorf("invoice creation failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Verificar que haya info en la respuesta
	if invoiceResp.Info == nil {
		c.log.Warn(ctx).
			Str("message", invoiceResp.Message).
			Msg("Invoice response has no info")
		return result, fmt.Errorf("invoice response has no info: %s", invoiceResp.Message)
	}

	// Verificar rechazo DIAN: error explícito o status "Rechazado"
	// Softpymes reporta rechazo con docsFe.error != "" O docsFe.status == "Rechazado"
	if invoiceResp.Info.DocsFe != nil {
		if invoiceResp.Info.DocsFe.Error != "" {
			c.log.Error(ctx).
				Str("docsFe_error", invoiceResp.Info.DocsFe.Error).
				Str("docsFe_status", invoiceResp.Info.DocsFe.Status).
				Str("document_number", invoiceResp.Info.DocumentNumber).
				Msg("DIAN rejected invoice (error field)")
			return result, fmt.Errorf("DIAN rejection: %s", invoiceResp.Info.DocsFe.Error)
		}
		if invoiceResp.Info.DocsFe.Status == "Rechazado" {
			c.log.Error(ctx).
				Str("docsFe_status", invoiceResp.Info.DocsFe.Status).
				Str("docsFe_message", invoiceResp.Info.DocsFe.Message).
				Str("document_number", invoiceResp.Info.DocumentNumber).
				Msg("DIAN rejected invoice (status Rechazado)")
			return result, fmt.Errorf("DIAN rejection: %s", invoiceResp.Info.DocsFe.Message)
		}
	}

	c.log.Info(ctx).
		Str("document_number", invoiceResp.Info.DocumentNumber).
		Str("date", invoiceResp.Info.Date).
		Float64("total", invoiceResp.Info.Total).
		Str("message", invoiceResp.Message).
		Msg("Invoice created successfully in Softpymes")

	// Rellenar resultado tipado
	result.InvoiceNumber = invoiceResp.Info.DocumentNumber
	result.ExternalID = invoiceResp.Info.DocumentNumber
	result.IssuedAt = invoiceResp.Info.Date

	result.ProviderInfo = map[string]interface{}{
		"subtotal":    invoiceResp.Info.Subtotal,
		"discount":    invoiceResp.Info.Discount,
		"iva":         invoiceResp.Info.IVA,
		"withholding": invoiceResp.Info.Withholding,
		"total":       invoiceResp.Info.Total,
	}

	if invoiceResp.Info.DocsFe != nil {
		result.ProviderInfo["dian_status"] = invoiceResp.Info.DocsFe.Status   // "Aceptado", "Pendiente", etc.
		result.ProviderInfo["dian_message"] = invoiceResp.Info.DocsFe.Message
	}

	return result, nil
}

// mapCurrencyToSoftpymes convierte códigos ISO de moneda al formato de Softpymes
// Softpymes usa: "P" = Peso Colombiano, "D" = Dólar Americano
func mapCurrencyToSoftpymes(isoCurrency string) string {
	switch isoCurrency {
	case "COP", "cop":
		return "P"
	case "USD", "usd":
		return "D"
	default:
		return "P"
	}
}

// findExistingInvoiceByOrderID busca en Softpymes una factura ya creada para una orden.
// La búsqueda se basa en el campo "comment" que almacena "order:<orderID>".
// Consulta los últimos 30 días (límite de la API de Softpymes).
// Retorna nil, nil si no se encuentra ninguna factura previa.
func (c *Client) findExistingInvoiceByOrderID(ctx context.Context, apiKey, apiSecret, referer, orderID, baseURL string) (*Document, error) {
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)
	dateTo := now.Format("2006-01-02")
	dateFrom := now.AddDate(0, 0, -30).Format("2006-01-02")

	params := ListDocumentsParams{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	docs, err := c.ListDocuments(ctx, apiKey, apiSecret, referer, params, baseURL)
	if err != nil {
		return nil, fmt.Errorf("error querying existing invoices: %w", err)
	}

	searchComment := "order:" + orderID
	for i, doc := range *docs {
		if strings.Contains(doc.Comment, searchComment) {
			c.log.Info(ctx).
				Str("order_id", orderID).
				Str("document_number", doc.DocumentNumber).
				Str("comment", doc.Comment).
				Msg("Found existing Softpymes invoice for order")
			return &(*docs)[i], nil
		}
	}

	return nil, nil
}
