package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	spDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ═══════════════════════════════════════════════════════════════
// DTOs locales replicados del módulo Invoicing para deserialización
// (Regla de aislamiento: no importar entre módulos)
// ═══════════════════════════════════════════════════════════════

// invoiceCustomerData datos del cliente (replicado de invoicing module)
type invoiceCustomerData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	Address string `json:"address,omitempty"`
}

// invoiceItemData datos de un item (replicado de invoicing module)
type invoiceItemData struct {
	ProductID   *string  `json:"product_id"`
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Quantity    int      `json:"quantity"`
	UnitPrice   float64  `json:"unit_price"`
	TotalPrice  float64  `json:"total_price"`
	Tax         float64  `json:"tax"`
	TaxRate     *float64 `json:"tax_rate"`
	Discount        float64  `json:"discount"`
	DiscountPercent float64  `json:"discount_percent"`
	// Precios en moneda presentment (moneda local, ej: COP)
	UnitPricePresentment  float64 `json:"unit_price_presentment"`
	TotalPricePresentment float64 `json:"total_price_presentment"`
	DiscountPresentment   float64 `json:"discount_presentment"`
	TaxPresentment        float64 `json:"tax_presentment"`
}

// invoiceData datos completos (replicado de invoicing module)
type invoiceData struct {
	IntegrationID uint                   `json:"integration_id"`
	Customer      invoiceCustomerData    `json:"customer"`
	Items         []invoiceItemData      `json:"items"`
	Total         float64                `json:"total"`
	Subtotal      float64                `json:"subtotal"`
	Tax           float64                `json:"tax"`
	Discount      float64                `json:"discount"`
	ShippingCost  float64                `json:"shipping_cost"`
	Currency      string                 `json:"currency"`
	OrderID       string                 `json:"order_id"`
	OrderNumber   string                 `json:"order_number,omitempty"`
	Config        map[string]interface{} `json:"config"`
}

// InvoiceRequestMessage es el mensaje recibido desde Invoicing Module
type InvoiceRequestMessage struct {
	InvoiceID     uint        `json:"invoice_id"`
	Provider      string      `json:"provider"`
	Operation     string      `json:"operation"`
	InvoiceData   invoiceData `json:"invoice_data"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"timestamp"`
}

// InvoiceRequestConsumer consume solicitudes de facturación desde Invoicing Module
type InvoiceRequestConsumer struct {
	rabbit            rabbitmq.IQueue
	integrationCore   integrationCore.IIntegrationService
	softpymesClient   ports.ISoftpymesClient
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del consumer
func New(
	rabbit rabbitmq.IQueue,
	integrationCore integrationCore.IIntegrationService,
	softpymesClient ports.ISoftpymesClient,
	responsePublisher *queue.ResponsePublisher,
	logger log.ILogger,
) *InvoiceRequestConsumer {
	return &InvoiceRequestConsumer{
		rabbit:            rabbit,
		integrationCore:   integrationCore,
		softpymesClient:   softpymesClient,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("softpymes.invoice_request_consumer"),
	}
}

const (
	QueueSoftpymesRequests = rabbitmq.QueueInvoicingSoftpymesRequests
)

// Start inicia el consumer
func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueSoftpymesRequests).
		Msg("Starting Softpymes invoice request consumer")

	// Declarar cola
	if err := c.rabbit.DeclareQueue(QueueSoftpymesRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	// Iniciar consumo
	if err := c.rabbit.Consume(ctx, QueueSoftpymesRequests, c.handleInvoiceRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	c.log.Info(ctx).
		Str("queue", QueueSoftpymesRequests).
		Msg("Consumer started successfully")

	return nil
}

// handleInvoiceRequest procesa una solicitud de facturación
func (c *InvoiceRequestConsumer) handleInvoiceRequest(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	// Parsear mensaje con DTOs tipados
	var request InvoiceRequestMessage
	if err := json.Unmarshal(message, &request); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("Failed to unmarshal request")
		return err
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("operation", request.Operation).
		Str("correlation_id", request.CorrelationID).
		Msg("Received invoice request")

	// Operación "compare": flujo propio (publica CompareResponseMessage, no InvoiceResponseMessage)
	if request.Operation == "compare" {
		return c.processCompareRequest(ctx, &request)
	}

	// Procesar según operación (create/retry/cancel/check_status → InvoiceResponseMessage)
	var response *queue.InvoiceResponseMessage
	switch request.Operation {
	case "create", "retry":
		response = c.processCreateInvoice(ctx, &request, startTime)
	case "check_status":
		response = c.processCheckStatus(ctx, &request, startTime)
	case "cancel":
		response = c.processCancelInvoice(ctx, &request, startTime)
	default:
		c.log.Warn(ctx).
			Str("operation", request.Operation).
			Msg("Unknown operation")
		response = c.createErrorResponse(&request, "unknown_operation", "Unknown operation: "+request.Operation, startTime, nil)
	}

	// Publicar response
	if err := c.responsePublisher.PublishResponse(ctx, response); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Failed to publish response")
		return err
	}

	return nil
}

// processCompareRequest obtiene documentos del proveedor en el rango de fechas
// y publica un CompareResponseMessage con todos los documentos encontrados.
func (c *InvoiceRequestConsumer) processCompareRequest(
	ctx context.Context,
	request *InvoiceRequestMessage,
) error {
	// 1. Extraer parámetros del Config
	dateFrom, _ := request.InvoiceData.Config["date_from"].(string)
	dateTo, _ := request.InvoiceData.Config["date_to"].(string)
	businessID := uint(0)
	if bid, ok := request.InvoiceData.Config["business_id"].(float64); ok {
		businessID = uint(bid)
	}

	c.log.Info(ctx).
		Str("date_from", dateFrom).
		Str("date_to", dateTo).
		Uint("business_id", businessID).
		Str("correlation_id", request.CorrelationID).
		Msg("Starting compare request")

	// Helper para publicar error en el canal de comparación
	publishErr := func(errMsg string) error {
		return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
			Operation:     "compare",
			CorrelationID: request.CorrelationID,
			BusinessID:    businessID,
			DateFrom:      dateFrom,
			DateTo:        dateTo,
			Error:         errMsg,
			Timestamp:     time.Now(),
		})
	}

	if dateFrom == "" || dateTo == "" {
		c.log.Error(ctx).Msg("date_from or date_to missing in compare config")
		return publishErr("date_from and date_to are required in compare config")
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in compare request")
		return publishErr("integration_id is 0")
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for compare")
		return publishErr("failed to get integration: " + err.Error())
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return publishErr("failed to decrypt api_key")
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return publishErr("failed to decrypt api_secret")
	}

	// 3. Combinar config de integración con config del mensaje
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	// 4. Resolver URL efectiva desde integration_type (base_url / base_url_test)
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integración Softpymes")
		return publishErr("base_url no configurada en el tipo de integración Softpymes (integration_types.base_url)")
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL for compare")

	// 5. Paginación: obtener todos los documentos del proveedor
	allDocs := make([]queue.CompareDocument, 0)
	pageSize := 20
	pageSizeStr := strconv.Itoa(pageSize)

	for page := 1; ; page++ {
		pageStr := strconv.Itoa(page)

		c.log.Info(ctx).
			Int("page", page).
			Str("date_from", dateFrom).
			Str("date_to", dateTo).
			Msg("Fetching documents page from Softpymes")

		docs, err := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
			DateFrom: dateFrom,
			DateTo:   dateTo,
			Page:     &pageStr,
			PageSize: &pageSizeStr,
		}, effectiveURL)
		if err != nil {
			c.log.Error(ctx).Err(err).Int("page", page).Msg("Failed to list documents")
			return publishErr(fmt.Sprintf("failed to list documents (page %d): %s", page, err.Error()))
		}

		for _, doc := range docs {
			details := make([]queue.CompareDocumentDetail, 0, len(doc.Details))
			for _, d := range doc.Details {
				details = append(details, queue.CompareDocumentDetail{
					ItemCode: d.ItemCode,
					ItemName: d.ItemName,
					Quantity: d.Quantity,
					Value:    d.Value,
					IVA:      d.IVA,
				})
			}
			allDocs = append(allDocs, queue.CompareDocument{
				DocumentNumber: doc.DocumentNumber,
				DocumentDate:   doc.DocumentDate,
				Total:          doc.Total,
				CustomerNit:    doc.CustomerNit,
				CustomerName:   doc.CustomerName,
				Comment:        doc.Comment,
				Prefix:         doc.Prefix,
				Details:        details,
			})
		}

		c.log.Info(ctx).
			Int("page", page).
			Int("page_count", len(docs)).
			Int("total_accumulated", len(allDocs)).
			Msg("Documents page fetched")

		// Última página cuando se devuelven menos registros que el tamaño de página
		if len(docs) < pageSize {
			break
		}
	}

	c.log.Info(ctx).
		Int("total_documents", len(allDocs)).
		Str("correlation_id", request.CorrelationID).
		Msg("All provider documents fetched, publishing compare response")

	// 6. Publicar resultado
	return c.responsePublisher.PublishCompareResponse(ctx, &queue.CompareResponseMessage{
		Operation:         "compare",
		CorrelationID:     request.CorrelationID,
		BusinessID:        businessID,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
		ProviderDocuments: allDocs,
		Timestamp:         time.Now(),
	})
}

// processCreateInvoice procesa la creación de una factura
func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Obtener integration_id directamente del DTO tipado
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in invoice_data")
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	// 2. Obtener integración desde IntegrationCore
	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration")
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime, nil)
	}

	// 3. Desencriptar credenciales
	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime, nil)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime, nil)
	}

	// 4. Combinar config de integración con config de facturación
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}

	// Sobrescribir con config específico de facturación
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	// 5. Resolver URL efectiva desde integration_type (base_url / base_url_test)
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).
			Uint("integration_id", integrationID).
			Msg("base_url no configurada en el tipo de integración Softpymes")
		return c.createErrorResponse(request, "missing_base_url",
			"base_url no configurada en el tipo de integración Softpymes (integration_types.base_url)",
			startTime, nil)
	}

	c.log.Info(ctx).
		Bool("is_testing", integration.IsTesting).
		Str("effective_url", effectiveURL).
		Msg("Resolved effective Softpymes URL")

	// 6. Construir request tipado para el cliente Softpymes
	invoiceReq := &spDtos.CreateInvoiceRequest{
		Customer: spDtos.CustomerData{
			Name:    request.InvoiceData.Customer.Name,
			Email:   request.InvoiceData.Customer.Email,
			Phone:   request.InvoiceData.Customer.Phone,
			DNI:     request.InvoiceData.Customer.DNI,
			Address: request.InvoiceData.Customer.Address,
		},
		Items:        mapItemsToClientDTOs(request.InvoiceData.Items),
		Total:        request.InvoiceData.Total,
		Subtotal:     request.InvoiceData.Subtotal,
		Tax:          request.InvoiceData.Tax,
		Discount:     request.InvoiceData.Discount,
		ShippingCost: request.InvoiceData.ShippingCost,
		Currency:     request.InvoiceData.Currency,
		OrderID:      request.InvoiceData.OrderID,
		OrderNumber:  request.InvoiceData.OrderNumber,
		Credentials: spDtos.Credentials{
			APIKey:    apiKey,
			APISecret: apiSecret,
		},
		Config:  combinedConfig,
		IsRetry: request.Operation == "retry",
	}

	// 7. Llamar al cliente HTTP de Softpymes con URL efectiva
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("effective_url", effectiveURL).
		Msg("Calling Softpymes API")

	result, err := c.softpymesClient.CreateInvoice(ctx, invoiceReq, effectiveURL)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Softpymes API call failed")

		// Incluir audit data del resultado (puede ser no-nil incluso en error)
		var auditData *spDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createErrorResponse(request, "api_error", err.Error(), startTime, auditData)
	}

	// 7b. Si Softpymes aceptó pero DIAN está validando, retornar como pending_validation
	if result.PendingValidation {
		c.log.Info(ctx).
			Uint("invoice_id", request.InvoiceID).
			Str("message", result.ProviderMessage).
			Msg("Invoice pending DIAN validation")

		processingTime := time.Since(startTime).Milliseconds()
		resp := &queue.InvoiceResponseMessage{
			InvoiceID:      request.InvoiceID,
			Provider:       "softpymes",
			Status:         "pending_validation",
			CorrelationID:  request.CorrelationID,
			Timestamp:      time.Now(),
			ProcessingTime: processingTime,
			Error:          result.ProviderMessage,
		}
		if result.AuditData != nil {
			resp.AuditRequestURL = result.AuditData.RequestURL
			resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
			resp.AuditResponseStatus = result.AuditData.ResponseStatus
			resp.AuditResponseBody = result.AuditData.ResponseBody
		}
		return resp
	}

	// 8. Consultar documento completo (GetDocumentByNumber)
	var fullDocument map[string]interface{}
	if result.InvoiceNumber != "" {
		referer, _ := combinedConfig["referer"].(string)

		c.log.Info(ctx).
			Str("invoice_number", result.InvoiceNumber).
			Msg("Waiting 3 seconds for DIAN processing")
		time.Sleep(3 * time.Second)

		c.log.Info(ctx).
			Str("invoice_number", result.InvoiceNumber).
			Msg("Fetching full document from Softpymes")

		doc, err := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, result.InvoiceNumber, effectiveURL)
		if err != nil {
			c.log.Warn(ctx).
				Err(err).
				Str("invoice_number", result.InvoiceNumber).
				Msg("Failed to fetch full document - continuing without it")
		} else {
			fullDocument = doc
			c.log.Info(ctx).
				Str("invoice_number", result.InvoiceNumber).
				Msg("Full document retrieved")

			// Usar el documentNumber del documento completo como número canónico.
			// La creación retorna el formato corto (ej: "FEV1001") pero el listado
			// del proveedor usa el formato padded (ej: "0000001001"). Guardamos el
			// padded para que la auditoría comparativa pueda cruzar ambos correctamente.
			if docNum, ok := fullDocument["documentNumber"].(string); ok && docNum != "" {
				c.log.Info(ctx).
					Str("old_invoice_number", result.InvoiceNumber).
					Str("new_invoice_number", docNum).
					Msg("Overriding invoice number with canonical padded format from full document")
				result.InvoiceNumber = docNum
				result.ExternalID = docNum
			}
		}
	}

	// 9. Send cash receipt if configured (non-fatal)
	referer, _ := combinedConfig["referer"].(string)
	c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

	// 10. Parsear issued_at
	var issuedAt *time.Time
	if result.IssuedAt != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, result.IssuedAt); parseErr == nil {
			issuedAt = &parsed
		}
	}

	// 11. Construir response exitosa con audit data
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		InvoiceNumber:  result.InvoiceNumber,
		ExternalID:     result.ExternalID,
		IssuedAt:       issuedAt,
		DocumentJSON:   fullDocument,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	// Incluir audit data en la respuesta
	if result.AuditData != nil {
		resp.AuditRequestURL = result.AuditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
		resp.AuditResponseStatus = result.AuditData.ResponseStatus
		resp.AuditResponseBody = result.AuditData.ResponseBody
	}

	return resp
}

// processCheckStatus busca un documento existente en Softpymes para una factura pendiente de DIAN.
// NO crea documentos nuevos — solo consulta ListDocuments y busca por comment "order:<UUID>".
// Si encuentra el documento, retorna success con el número de factura.
// Si no lo encuentra, retorna pending_validation para programar otro check más tarde.
func (c *InvoiceRequestConsumer) processCheckStatus(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	orderID := request.InvoiceData.OrderID
	if orderID == "" {
		return c.createErrorResponse(request, "missing_order_id", "order_id is required for check_status", startTime, nil)
	}

	// 1. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime, nil)
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime, nil)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime, nil)
	}

	// 2. Resolver config y URL
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		return c.createErrorResponse(request, "missing_base_url", "base_url no configurada", startTime, nil)
	}

	// 3. Buscar documento por order_id en comment field
	branchCode := "001"
	if bc, ok := combinedConfig["branch_code"].(string); ok && bc != "" {
		branchCode = bc
	}

	// Rango: desde hace 30 días hasta hoy (máximo permitido por Softpymes)
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().In(loc)
	dateTo := now.Format("2006-01-02")
	dateFrom := now.AddDate(0, 0, -30).Format("2006-01-02")

	pageSize := "50"
	docs, err := c.softpymesClient.ListDocuments(ctx, apiKey, apiSecret, referer, ports.ListDocumentsParams{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		PageSize: &pageSize,
	}, effectiveURL)

	if err != nil {
		c.log.Warn(ctx).Err(err).
			Uint("invoice_id", request.InvoiceID).
			Str("order_id", orderID).
			Msg("Failed to search documents in Softpymes — keeping as pending")

		// No marcamos como error — es una falla temporal de conectividad
		processingTime := time.Since(startTime).Milliseconds()
		return &queue.InvoiceResponseMessage{
			InvoiceID:      request.InvoiceID,
			Provider:       "softpymes",
			Status:         "pending_validation",
			CorrelationID:  request.CorrelationID,
			Timestamp:      time.Now(),
			ProcessingTime: processingTime,
			Error:          "Check status failed, will retry: " + err.Error(),
		}
	}

	// 4. Buscar documento con comment "order:<UUID>"
	searchComment := "order:" + orderID
	_ = branchCode // branchCode ya está en el filtro implícito de la búsqueda
	for _, doc := range docs {
		if strings.Contains(doc.Comment, searchComment) {
			c.log.Info(ctx).
				Uint("invoice_id", request.InvoiceID).
				Str("order_id", orderID).
				Str("document_number", doc.DocumentNumber).
				Msg("Found existing document in Softpymes for pending invoice")

			// 5. Obtener documento completo
			var fullDocument map[string]interface{}
			if doc.DocumentNumber != "" {
				fullDoc, err := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, doc.DocumentNumber, effectiveURL)
				if err != nil {
					c.log.Warn(ctx).Err(err).Str("document_number", doc.DocumentNumber).Msg("Failed to get full document")
				} else {
					fullDocument = fullDoc
					// Usar documentNumber canónico del full document
					if docNum, ok := fullDocument["documentNumber"].(string); ok && docNum != "" {
						doc.DocumentNumber = docNum
					}
				}
			}

			// 6. Send cash receipt if configured (non-fatal)
			c.sendCashReceiptIfConfigured(ctx, fullDocument, combinedConfig, apiKey, apiSecret, referer, effectiveURL, request.InvoiceID)

			// 7. Parsear fecha
			var issuedAt *time.Time
			if doc.DocumentDate != "" {
				if parsed, parseErr := time.Parse("2006-01-02", doc.DocumentDate); parseErr == nil {
					issuedAt = &parsed
				}
			}

			processingTime := time.Since(startTime).Milliseconds()
			return &queue.InvoiceResponseMessage{
				InvoiceID:      request.InvoiceID,
				Provider:       "softpymes",
				Status:         "success",
				InvoiceNumber:  doc.DocumentNumber,
				ExternalID:     doc.DocumentNumber,
				IssuedAt:       issuedAt,
				DocumentJSON:   fullDocument,
				CorrelationID:  request.CorrelationID,
				Timestamp:      time.Now(),
				ProcessingTime: processingTime,
			}
		}
	}

	// 7. No encontrado — DIAN sigue validando, mantener pending
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", orderID).
		Int("docs_searched", len(docs)).
		Msg("Document not found in Softpymes — DIAN still validating, keeping pending")

	processingTime := time.Since(startTime).Milliseconds()
	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "pending_validation",
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
		Error:          "Document not found yet — DIAN still validating",
	}
}

// processCancelInvoice anula una factura emitida en Softpymes
func (c *InvoiceRequestConsumer) processCancelInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Extraer external_id y reason del Config
	documentNumber, _ := request.InvoiceData.Config["external_id"].(string)
	if documentNumber == "" {
		c.log.Error(ctx).Uint("invoice_id", request.InvoiceID).Msg("external_id missing in cancel config")
		return c.createCancelErrorResponse(request, "missing_external_id", "external_id is required for cancellation", startTime)
	}

	reason, _ := request.InvoiceData.Config["cancel_reason"].(string)
	if reason == "" {
		reason = "Anulación de factura"
	}

	// 2. Obtener integración y credenciales
	integrationID := request.InvoiceData.IntegrationID
	if integrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in cancel request")
		return c.createCancelErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime)
	}

	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration for cancel")
		return c.createCancelErrorResponse(request, "integration_not_found", err.Error(), startTime)
	}

	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return c.createCancelErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return c.createCancelErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime)
	}

	// 3. Combinar config
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}
	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	referer, _ := combinedConfig["referer"].(string)

	// 4. Resolver URL efectiva
	effectiveURL := integration.BaseURL
	if integration.IsTesting && integration.BaseURLTest != "" {
		effectiveURL = integration.BaseURLTest
	}
	if effectiveURL == "" {
		c.log.Error(ctx).Uint("integration_id", integrationID).Msg("base_url no configurada para cancelación")
		return c.createCancelErrorResponse(request, "missing_base_url", "base_url no configurada en el tipo de integración Softpymes", startTime)
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("document_number", documentNumber).
		Str("effective_url", effectiveURL).
		Msg("Cancelling invoice in Softpymes")

	// 5. Llamar al cliente Softpymes
	if err := c.softpymesClient.CancelInvoice(ctx, apiKey, apiSecret, referer, documentNumber, reason, effectiveURL); err != nil {
		c.log.Error(ctx).Err(err).Uint("invoice_id", request.InvoiceID).Msg("Softpymes cancel failed")
		return c.createCancelErrorResponse(request, "api_error", err.Error(), startTime)
	}

	processingTime := time.Since(startTime).Milliseconds()
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("document_number", documentNumber).
		Msg("Invoice cancelled successfully in Softpymes")

	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		Operation:      "cancel",
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}
}

// createCancelErrorResponse crea una respuesta de error para operaciones de cancelación
func (c *InvoiceRequestConsumer) createCancelErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "error",
		Operation:      "cancel",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: time.Since(startTime).Milliseconds(),
	}
}

// mapItemsToClientDTOs convierte items del mensaje a DTOs del cliente Softpymes
func mapItemsToClientDTOs(items []invoiceItemData) []spDtos.ItemData {
	result := make([]spDtos.ItemData, 0, len(items))
	for _, item := range items {
		result = append(result, spDtos.ItemData{
			ProductID:             item.ProductID,
			SKU:                   item.SKU,
			Name:                  item.Name,
			Description:           item.Description,
			Quantity:              item.Quantity,
			UnitPrice:             item.UnitPrice,
			TotalPrice:            item.TotalPrice,
			Tax:                   item.Tax,
			TaxRate:               item.TaxRate,
			Discount:              item.Discount,
			DiscountPercent:       item.DiscountPercent,
			UnitPricePresentment:  item.UnitPricePresentment,
			TotalPricePresentment: item.TotalPricePresentment,
			DiscountPresentment:   item.DiscountPresentment,
			TaxPresentment:        item.TaxPresentment,
		})
	}
	return result
}

// createErrorResponse crea una respuesta de error, opcionalmente con audit data
func (c *InvoiceRequestConsumer) createErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
	auditData *spDtos.AuditData,
) *queue.InvoiceResponseMessage {
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "error",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	// Incluir audit data si está disponible (ej: cuando el HTTP request se hizo pero falló)
	if auditData != nil {
		resp.AuditRequestURL = auditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(auditData.RequestPayload)
		resp.AuditResponseStatus = auditData.ResponseStatus
		resp.AuditResponseBody = auditData.ResponseBody
	}

	return resp
}

// sendCashReceiptIfConfigured envía un recibo de caja si la config lo tiene habilitado.
// Es non-fatal: si falla, se loguea el error pero no afecta el resultado de la factura.
func (c *InvoiceRequestConsumer) sendCashReceiptIfConfigured(
	ctx context.Context,
	fullDocument map[string]interface{},
	config map[string]interface{},
	apiKey, apiSecret, referer, baseURL string,
	invoiceID uint,
) {
	sendCashReceipt, _ := config["send_cash_receipt"].(bool)
	if !sendCashReceipt {
		return
	}

	if fullDocument == nil {
		c.log.Warn(ctx).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt configured but full document is nil — skipping")
		return
	}

	c.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Msg("Sending cash receipt (configured in integration)")

	if err := c.softpymesClient.SendCashReceiptFromDocument(ctx, apiKey, apiSecret, referer, baseURL, fullDocument, config); err != nil {
		c.log.Error(ctx).Err(err).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt failed — invoice created but payment not registered in Softpymes")
	} else {
		c.log.Info(ctx).
			Uint("invoice_id", invoiceID).
			Msg("Cash receipt sent successfully")
	}
}

// toMapPayload convierte cualquier valor (struct o map) a map[string]interface{} via JSON.
func toMapPayload(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
