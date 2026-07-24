package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
)

type InvoiceResponse struct {
	Message string       `json:"message"`
	Info    *InvoiceInfo `json:"info,omitempty"`
}

type InvoiceInfo struct {
	Date           string  `json:"date"`
	DocumentNumber string  `json:"documentNumber"`
	Subtotal       float64 `json:"subtotal"`
	Discount       float64 `json:"discount"`
	IVA            float64 `json:"iva"`
	Withholding    float64 `json:"withholding"`
	Total          float64 `json:"total"`
	DocsFe         *DocsFe `json:"docsFe,omitempty"`
}

type DocsFe struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (c *Client) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest, baseURL string) (*dtos.CreateInvoiceResult, error) {
	result := &dtos.CreateInvoiceResult{}

	if req.Credentials.APIKey == "" {
		return result, fmt.Errorf("api_key not found in credentials")
	}
	if req.Credentials.APISecret == "" {
		return result, fmt.Errorf("api_secret not found in credentials")
	}

	referer, _ := req.Config["referer"].(string)
	if referer == "" {
		return result, fmt.Errorf("referer not found in config")
	}

	token, err := c.authenticate(ctx, req.Credentials.APIKey, req.Credentials.APISecret, referer, baseURL)
	if err != nil {
		return result, fmt.Errorf("authentication failed: %w", err)
	}

	if req.IsRetry && req.OrderID != "" {
		branchCode := "001"
		if bc, ok := req.Config["branch_code"].(string); ok && bc != "" {
			branchCode = bc
		}
		existing, err := c.findExistingInvoiceByOrderID(ctx, req.Credentials.APIKey, req.Credentials.APISecret, referer, req.OrderID, branchCode, req.OrderCreatedAt, baseURL)
		if err != nil {
			c.log.Error(ctx).Err(err).
				Str("order_id", req.OrderID).
				Msg("Could not verify existing invoice in Softpymes, aborting retry to avoid duplicates")
			return result, fmt.Errorf("no se pudo verificar si la factura ya existe en Softpymes, reintento abortado para evitar duplicados: %w", err)
		}
		if existing != nil {
			c.log.Info(ctx).
				Str("order_id", req.OrderID).
				Str("document_number", existing.DocumentNumber).
				Msg("Invoice already exists for this order in Softpymes, skipping duplicate creation")
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

	forceDefault, _ := req.Config["force_default_customer"].(bool)
	customerNit := ""
	if forceDefault {
		if defaultNit, ok := req.Config["default_customer_nit"].(string); ok && defaultNit != "" {
			customerNit = defaultNit
			c.log.Info(ctx).
				Str("default_nit", defaultNit).
				Str("original_dni", req.Customer.DNI).
				Msg("Using forced default customer NIT (force_default_customer=true)")
		} else {
			c.log.Error(ctx).
				Msg("force_default_customer is true but default_customer_nit is not configured")
			return result, fmt.Errorf("force_default_customer is true but default_customer_nit is not configured")
		}
	} else {
		customerNit = req.Customer.DNI
		if customerNit == "" {
			if defaultNit, ok := req.Config["default_customer_nit"].(string); ok && defaultNit != "" {
				customerNit = defaultNit
				c.log.Info(ctx).
					Str("default_nit", defaultNit).
					Msg("Using default customer NIT from config (customer has no DNI)")
			} else {
				c.log.Error(ctx).
					Msg("customerNit is required but not provided. Configure default_customer_nit in integration config or ensure customers have DNI")
				return result, fmt.Errorf("customerNit is required: customer has no DNI and no default_customer_nit configured")
			}
		}
	}

	customerBranch := ""
	if branch, err := c.ensureCustomerExists(ctx, token, referer, customerNit, &req.Customer, req.Config, baseURL); err != nil {
		c.log.Warn(ctx).Err(err).
			Str("customer_nit", customerNit).
			Msg("Could not ensure customer exists in Softpymes, proceeding anyway")
	} else {
		customerBranch = branch
	}

	if customerBranch == "" {
		customerBranch = "001"
		if cb, ok := req.Config["customer_branch_code"].(string); ok && cb != "" {
			customerBranch = cb
		}
	}

	branchCode := "001"
	if branch, ok := req.Config["branch_code"].(string); ok && branch != "" {
		branchCode = branch
	}

	sellerNit := ""
	if seller, ok := req.Config["seller_nit"].(string); ok && seller != "" {
		sellerNit = seller
	}

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

	if len(req.Items) == 0 {
		c.log.Error(ctx).Msg("No items found in invoice request - cannot create invoice without items")
		return result, fmt.Errorf("no items found in invoice data")
	}

	var itemMappings map[string]interface{}
	if mappings, ok := req.Config["item_mappings"].(map[string]interface{}); ok {
		itemMappings = mappings
	}

	softpymesItems := make([]map[string]interface{}, 0, len(req.Items))
	for _, item := range req.Items {
		itemCode := resolveItemCode(item.SKU, item.Name, item.ProductID, itemMappings)

		unitPrice := item.UnitPriceBase

		if unitPrice == 0 && item.UnitPrice > 0 {
			rate := 0.19
			if item.TaxRate != nil && *item.TaxRate > 0 {
				rate = *item.TaxRate
			}
			unitPrice = item.UnitPrice / (1 + rate)
		}

		if unitPrice == 0 {
			unitPrice = item.UnitPrice
		}

		softpymesItem := map[string]interface{}{
			"itemCode":  itemCode,
			"quantity":  float64(item.Quantity),
			"unitCode":  "UNI",
			"discount":  item.DiscountPercent,
			"unitValue": fmt.Sprintf("%.2f", unitPrice),
		}

		softpymesItems = append(softpymesItems, softpymesItem)
	}

	effectiveShipping := req.ShippingCost - req.ShippingDiscount
	if effectiveShipping < 0 {
		effectiveShipping = 0
	}

	if effectiveShipping > 0 {
		shippingPrice := req.ShippingCostBase

		if shippingPrice == 0 && effectiveShipping > 0 {
			shippingPrice = effectiveShipping / 1.19
		}

		if shippingPrice == 0 {
			shippingPrice = effectiveShipping
		}

		shippingItem := map[string]interface{}{
			"itemCode":  resolveShippingItemCode(itemMappings),
			"quantity":  1.0,
			"unitCode":  "UNI",
			"discount":  0,
			"unitValue": fmt.Sprintf("%.2f", shippingPrice),
		}
		softpymesItems = append(softpymesItems, shippingItem)
	}

	rawCurrency := req.Currency
	if rawCurrency == "" {
		rawCurrency = "COP"
	}
	currency := mapCurrencyToSoftpymes(rawCurrency)

	loc, _ := time.LoadLocation("America/Bogota")
	documentDate := time.Now().In(loc).Format("2006-01-02")

	comment := ""
	if req.OrderID != "" {
		comment = "order:" + req.OrderID
		if req.OrderNumber != "" {
			comment += " | #" + req.OrderNumber
		}
	}

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
		Post(requestURL)

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

	if invoiceResp.Info == nil {
		if strings.Contains(strings.ToLower(invoiceResp.Message), "validación") {
			c.log.Info(ctx).
				Str("message", invoiceResp.Message).
				Msg("Invoice accepted by Softpymes, pending DIAN validation")
			result.PendingValidation = true
			result.ProviderMessage = invoiceResp.Message
			return result, nil
		}
		c.log.Warn(ctx).
			Str("message", invoiceResp.Message).
			Msg("Invoice response has no info")
		return result, fmt.Errorf("invoice response has no info: %s", invoiceResp.Message)
	}

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
		result.ProviderInfo["dian_status"] = invoiceResp.Info.DocsFe.Status
		result.ProviderInfo["dian_message"] = invoiceResp.Info.DocsFe.Message
	}

	return result, nil
}

func (c *Client) CancelInvoice(ctx context.Context, apiKey, apiSecret, referer, documentNumber, reason, baseURL string) error {
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	cancelReq := map[string]interface{}{
		"documentNumber": documentNumber,
		"reason":         reason,
	}

	var cancelResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	requestURL := c.resolveURL(baseURL, "/app/integration/sales_invoice/cancel/")
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetBody(cancelReq).
		SetResult(&cancelResp).
		Post(requestURL)

	if err != nil {
		c.log.Error(ctx).Err(err).Str("document_number", documentNumber).Msg("Failed to cancel invoice")
		return fmt.Errorf("invoice cancellation request failed: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return fmt.Errorf("authentication token expired")
		}
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("Invoice cancellation failed")
		return fmt.Errorf("invoice cancellation failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	if cancelResp.Error != "" {
		return fmt.Errorf("invoice cancellation error: %s", cancelResp.Error)
	}

	c.log.Info(ctx).
		Str("document_number", documentNumber).
		Str("message", cancelResp.Message).
		Msg("Invoice cancelled successfully in Softpymes")

	return nil
}

func resolveItemCode(itemSKU string, itemName string, productID *string, itemMappings map[string]interface{}) string {
	if itemMappings != nil && itemName != "" {
		nameLower := strings.ToLower(strings.TrimSpace(itemName))
		if nameLower == "tip" || nameLower == "propina" {
			if code, ok := itemMappings["tip"].(string); ok && code != "" {
				return code
			}
		}
		if strings.Contains(nameLower, "membership") || strings.Contains(nameLower, "membresía") || strings.Contains(nameLower, "membresia") {
			if code, ok := itemMappings["membership"].(string); ok && code != "" {
				return code
			}
		}
	}
	if itemSKU != "" {
		return itemSKU
	}
	if productID != nil {
		return *productID
	}
	return ""
}

func resolveShippingItemCode(itemMappings map[string]interface{}) string {
	if itemMappings != nil {
		if code, ok := itemMappings["shipping"].(string); ok && code != "" {
			return code
		}
	}
	return "SHIPPING"
}

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

const (
	existingInvoiceSearchPageSize = 50
	existingInvoiceSearchMaxPages = 100
)

func (c *Client) findExistingInvoiceByOrderID(ctx context.Context, apiKey, apiSecret, referer, orderID, branchCode, orderCreatedAt, baseURL string) (*Document, error) {
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)
	dateTo := now.Format("2006-01-02")

	dateFrom := dateTo
	if orderCreatedAt != "" {
		dateFrom = orderCreatedAt
	}

	searchComment := "order:" + orderID
	pageSize := strconv.Itoa(existingInvoiceSearchPageSize)

	c.log.Info(ctx).
		Str("order_id", orderID).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Str("branch_code", branchCode).
		Msg("Checking for existing Softpymes invoice (retry idempotency)")

	for page := 1; page <= existingInvoiceSearchMaxPages; page++ {
		pageStr := strconv.Itoa(page)
		params := ListDocumentsParams{
			DateFrom: dateFrom,
			DateTo:   dateTo,
			PageSize: &pageSize,
			Page:     &pageStr,
		}
		if branchCode != "" {
			params.BranchCode = &branchCode
		}

		docs, err := c.listDocuments(ctx, apiKey, apiSecret, referer, params, baseURL)
		if err != nil {
			return nil, fmt.Errorf("error querying existing invoices (page %d): %w", page, err)
		}
		if docs == nil {
			return nil, nil
		}

		for i, doc := range *docs {
			if strings.Contains(doc.Comment, searchComment) {
				c.log.Info(ctx).
					Str("order_id", orderID).
					Str("document_number", doc.DocumentNumber).
					Int("page", page).
					Msg("Found existing Softpymes invoice for order — skipping duplicate")
				return &(*docs)[i], nil
			}
		}

		if len(*docs) < existingInvoiceSearchPageSize {
			return nil, nil
		}
	}

	c.log.Warn(ctx).
		Str("order_id", orderID).
		Int("max_pages", existingInvoiceSearchMaxPages).
		Msg("Existing invoice search exhausted max pages without a conclusive result")
	return nil, fmt.Errorf("no fue posible verificar de forma concluyente si la orden %s ya tiene factura (se agotaron %d paginas de busqueda)", orderID, existingInvoiceSearchMaxPages)
}
