package client

import (
	"context"
	"fmt"
	"time"
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
type DocsFe struct {
	Status  bool   `json:"status"`  // true = válido
	Message string `json:"message"` // "Documento válido enviado al proveedor tecnológico"
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

	// Preparar request de factura según documentación oficial de Softpymes
	// https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostSaleInvoice

	// Extraer datos del cliente
	customer, _ := invoiceData["customer"].(map[string]interface{})
	customerNit := ""
	if customer != nil {
		if dni, ok := customer["dni"].(string); ok && dni != "" {
			customerNit = dni
		}
	}

	// Si el customer no tiene NIT, intentar obtener uno por defecto del config
	if customerNit == "" {
		if defaultNit, ok := config["default_customer_nit"].(string); ok && defaultNit != "" {
			customerNit = defaultNit
			c.log.Info(ctx).
				Str("default_nit", defaultNit).
				Msg("Using default customer NIT from config")
		} else {
			// Si no hay NIT ni en el customer ni en el config, log error y retornar
			c.log.Error(ctx).
				Msg("❌ customerNit is required but not provided. Configure default_customer_nit in integration config or ensure customers have DNI")
			return fmt.Errorf("customerNit is required: customer has no DNI and no default_customer_nit configured")
		}
	}

	// Asegurar que el cliente existe en Softpymes antes de facturar
	// Si no existe, crearlo automáticamente usando la API de Softpymes
	if err := c.ensureCustomerExists(ctx, token, referer, customerNit, customer, config); err != nil {
		c.log.Warn(ctx).Err(err).
			Str("customer_nit", customerNit).
			Msg("⚠️ Could not ensure customer exists in Softpymes, proceeding anyway")
	}

	// Obtener branch_code del config (default "001" si no existe)
	branchCode := "001" // Default
	if branch, ok := config["branch_code"].(string); ok && branch != "" {
		branchCode = branch
	}

	// Obtener customer_branch_code del config (default "000" si no existe)
	// REQUERIDO por Softpymes: código de la sucursal del cliente
	// Softpymes asigna "000" por defecto al crear un cliente nuevo
	customerBranch := "000" // Default (Softpymes genera "000" al crear cliente)
	if cb, ok := config["customer_branch_code"].(string); ok && cb != "" {
		customerBranch = cb
	}

	// Obtener seller_nit del config (OPCIONAL - solo enviar si está configurado)
	// Softpymes requiere que el seller exista en su sistema via /app/integration/seller
	// Si no hay sellers configurados, NO enviar el campo para evitar 404
	sellerNit := ""
	if seller, ok := config["seller_nit"].(string); ok && seller != "" {
		sellerNit = seller
	}

	// Obtener resolution_id del config
	// IMPORTANTE: resolutionId debe ser un ID válido obtenido desde /app/integration/resolutions
	// Si es 0, Softpymes puede rechazar la factura
	resolutionID := 0
	if resID, ok := config["resolution_id"].(float64); ok {
		resolutionID = int(resID)
	} else if resID, ok := config["resolution_id"].(int); ok {
		resolutionID = resID
	}

	// Log warning si resolution_id es 0
	if resolutionID == 0 {
		c.log.Warn(ctx).
			Msg("⚠️ resolutionId is 0 - Softpymes may reject this invoice. Configure a valid resolution_id in integration config")
	}

	// Mapear items al formato de Softpymes
	// Después de JSON marshal/unmarshal (RabbitMQ), []map[string]interface{} se convierte en []interface{}
	var items []map[string]interface{}
	if rawItems, ok := invoiceData["items"].([]interface{}); ok {
		for _, rawItem := range rawItems {
			if item, ok := rawItem.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
	} else if directItems, ok := invoiceData["items"].([]map[string]interface{}); ok {
		items = directItems
	}

	if len(items) == 0 {
		c.log.Error(ctx).
			Interface("items_raw", invoiceData["items"]).
			Msg("❌ No items found in invoice data - cannot create invoice without items")
		return fmt.Errorf("no items found in invoice data")
	}
	softpymesItems := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		// Extraer datos del item
		itemCode := ""
		if sku, ok := item["sku"].(string); ok {
			itemCode = sku
		} else if productID, ok := item["product_id"].(string); ok {
			itemCode = productID
		}

		quantity := 0.0
		if q, ok := item["quantity"].(int); ok {
			quantity = float64(q)
		} else if q, ok := item["quantity"].(float64); ok {
			quantity = q
		}

		discount := 0.0
		if d, ok := item["discount"].(float64); ok {
			discount = d
		} else if d, ok := item["discount"].(int); ok {
			discount = float64(d)
		}

		// unitCode por defecto "UNI" (UNIDADES) - código estándar en Softpymes
		// Verificado via API: /app/integration/items/:code → units[].code = "UNI"
		unitCode := "UNI"
		if unit, ok := item["unit_code"].(string); ok && unit != "" {
			unitCode = unit
		}

		// unitValue (precio unitario) - enviarlo para usar nuestro precio en vez del price list de Softpymes
		unitValue := 0.0
		if up, ok := item["unit_price"].(float64); ok {
			unitValue = up
		} else if up, ok := item["unit_price"].(int); ok {
			unitValue = float64(up)
		}

		softpymesItem := map[string]interface{}{
			"itemCode": itemCode,
			"quantity": quantity,
			"discount": discount,
			"unitCode": unitCode,
		}

		// Solo incluir unitValue si tiene valor (Softpymes lo usa para override del precio de lista)
		// IMPORTANTE: Softpymes espera unitValue como String, no como número
		if unitValue > 0 {
			softpymesItem["unitValue"] = fmt.Sprintf("%.2f", unitValue)
		}

		softpymesItems = append(softpymesItems, softpymesItem)
	}

	// Obtener currency y mapear al formato de Softpymes
	// Softpymes usa códigos propios: "P" = Peso Colombiano, "D" = Dólar Americano
	rawCurrency := "COP"
	if curr, ok := invoiceData["currency"].(string); ok && curr != "" {
		rawCurrency = curr
	}
	currency := mapCurrencyToSoftpymes(rawCurrency)

	// Generar documentDate en formato YYYY-MM-DD (zona horaria Colombia UTC-5)
	// REQUERIDO por Softpymes - sin este campo retorna 500
	loc, _ := time.LoadLocation("America/Bogota")
	documentDate := time.Now().In(loc).Format("2006-01-02")

	// Construir request según formato de Softpymes
	// Docs: https://api-integracion.softpymes.com.co/doc/#api-Documentos-PostSaleInvoice
	invoiceReq := map[string]interface{}{
		"documentDate":   documentDate, // REQUERIDO: formato YYYY-MM-DD, zona Colombia
		"currencyCode":   currency,
		"exchangeRate":   1.0, // Por ahora siempre 1.0 (moneda local)
		"branchCode":     branchCode,
		"customerBranch": customerBranch, // REQUERIDO: código sucursal del cliente
		"customerNit":    customerNit,
		"termDays":       0, // Por ahora siempre contado (0 días)
		"resolutionId":   resolutionID,
		"comment":        "", // Observaciones del documento
		"items":          softpymesItems,
	}

	// Solo incluir sellerNit si está configurado
	// Si no hay sellers en Softpymes, enviar el campo causa 404
	if sellerNit != "" {
		invoiceReq["sellerNit"] = sellerNit
	}

	// Log detallado del request para debugging
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
		Msg("📤 Sending invoice request to Softpymes")

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
			Str("response", string(resp.Body())).
			Msg("Invoice creation failed")

		// Si es 401, el token expiró
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return fmt.Errorf("authentication token expired")
		}

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

	return nil
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
		return "P" // Default a Pesos
	}
}
