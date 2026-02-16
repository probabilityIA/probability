package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// InvoiceRequestMessage es el mensaje recibido desde Invoicing Module
type InvoiceRequestMessage struct {
	InvoiceID     uint                   `json:"invoice_id"`
	Provider      string                 `json:"provider"`
	Operation     string                 `json:"operation"`
	InvoiceData   map[string]interface{} `json:"invoice_data"`
	CorrelationID string                 `json:"correlation_id"`
	Timestamp     time.Time              `json:"timestamp"`
}

// InvoiceRequestConsumer consume solicitudes de facturación desde Invoicing Module
type InvoiceRequestConsumer struct {
	rabbit            rabbitmq.IQueue
	integrationCore   integrationCore.IIntegrationCore
	softpymesClient   ports.ISoftpymesClient
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
}

// NewInvoiceRequestConsumer crea una nueva instancia del consumer
func NewInvoiceRequestConsumer(
	rabbit rabbitmq.IQueue,
	integrationCore integrationCore.IIntegrationCore,
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
	QueueSoftpymesRequests = "invoicing.softpymes.requests"
)

// Start inicia el consumer
func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueSoftpymesRequests).
		Msg("🚀 Starting Softpymes invoice request consumer")

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
		Msg("✅ Consumer started successfully")

	return nil
}

// handleInvoiceRequest procesa una solicitud de facturación
func (c *InvoiceRequestConsumer) handleInvoiceRequest(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	// Parsear mensaje
	var request InvoiceRequestMessage
	if err := json.Unmarshal(message, &request); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("❌ Failed to unmarshal request")
		return err
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("operation", request.Operation).
		Str("correlation_id", request.CorrelationID).
		Msg("📩 Received invoice request")

	// Procesar según operación
	var response *queue.InvoiceResponseMessage
	switch request.Operation {
	case "create", "retry":
		response = c.processCreateInvoice(ctx, &request, startTime)
	default:
		c.log.Warn(ctx).
			Str("operation", request.Operation).
			Msg("⚠️ Unknown operation")
		response = c.createErrorResponse(&request, "unknown_operation", "Unknown operation: "+request.Operation, startTime)
	}

	// Publicar response
	if err := c.responsePublisher.PublishResponse(ctx, response); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("❌ Failed to publish response")
		return err
	}

	return nil
}

// processCreateInvoice procesa la creación de una factura
func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Obtener integration_id del invoice_data
	integrationIDFloat, ok := request.InvoiceData["integration_id"].(float64)
	if !ok {
		c.log.Error(ctx).Msg("integration_id not found in invoice_data")
		return c.createErrorResponse(request, "missing_integration_id", "integration_id not found", startTime)
	}
	integrationID := uint(integrationIDFloat)

	// 2. Obtener integración desde IntegrationCore
	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := c.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get integration")
		return c.createErrorResponse(request, "integration_not_found", err.Error(), startTime)
	}

	// 3. Desencriptar credenciales
	apiKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_key", startTime)
	}

	apiSecret, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt api_secret", startTime)
	}

	// 4. Construir invoiceData completo con credentials y config
	invoiceData := request.InvoiceData

	// Agregar credentials
	invoiceData["credentials"] = map[string]interface{}{
		"api_key":    apiKey,
		"api_secret": apiSecret,
	}

	// Combinar config de integración con config de facturación
	combinedConfig := make(map[string]interface{})

	// Primero copiar config de integración (Softpymes: referer, api_url, company_nit)
	if integration.Config != nil {
		if configMap, ok := integration.Config.(map[string]interface{}); ok {
			for k, v := range configMap {
				combinedConfig[k] = v
			}
		}
	}

	// Luego sobrescribir con config específico de facturación
	if invoiceConfig, ok := invoiceData["config"].(map[string]interface{}); ok {
		for k, v := range invoiceConfig {
			combinedConfig[k] = v
		}
	}

	invoiceData["config"] = combinedConfig

	// 5. Llamar al cliente HTTP de Softpymes
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Msg("📡 Calling Softpymes API")

	if err := c.softpymesClient.CreateInvoice(ctx, invoiceData); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("❌ Softpymes API call failed")
		return c.createErrorResponse(request, "api_error", err.Error(), startTime)
	}

	// 6. Extraer datos de la respuesta (modificados in-place por el cliente)
	invoiceNumber, _ := invoiceData["invoice_number"].(string)
	externalID, _ := invoiceData["external_id"].(string)
	invoiceURL, _ := invoiceData["invoice_url"].(string)
	pdfURL, _ := invoiceData["pdf_url"].(string)
	xmlURL, _ := invoiceData["xml_url"].(string)
	cufe, _ := invoiceData["cufe"].(string)

	var issuedAt *time.Time
	if issuedAtStr, ok := invoiceData["issued_at"].(string); ok && issuedAtStr != "" {
		if parsed, err := time.Parse(time.RFC3339, issuedAtStr); err == nil {
			issuedAt = &parsed
		}
	}

	// 7. Consultar documento completo (GetDocumentByNumber)
	var fullDocument map[string]interface{}
	if invoiceNumber != "" {
		referer, _ := combinedConfig["referer"].(string)

		c.log.Info(ctx).
			Str("invoice_number", invoiceNumber).
			Msg("⏳ Waiting 3 seconds for DIAN processing")
		time.Sleep(3 * time.Second)

		c.log.Info(ctx).
			Str("invoice_number", invoiceNumber).
			Msg("📥 Fetching full document from Softpymes")

		doc, err := c.softpymesClient.GetDocumentByNumber(ctx, apiKey, apiSecret, referer, invoiceNumber)
		if err != nil {
			c.log.Warn(ctx).
				Err(err).
				Str("invoice_number", invoiceNumber).
				Msg("⚠️ Failed to fetch full document - continuing without it")
		} else {
			fullDocument = doc
			c.log.Info(ctx).
				Str("invoice_number", invoiceNumber).
				Msg("✅ Full document retrieved")
		}
	}

	// 8. Construir response exitosa
	processingTime := time.Since(startTime).Milliseconds()

	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "success",
		InvoiceNumber:  invoiceNumber,
		ExternalID:     externalID,
		InvoiceURL:     invoiceURL,
		PDFURL:         pdfURL,
		XMLURL:         xmlURL,
		CUFE:           cufe,
		IssuedAt:       issuedAt,
		DocumentJSON:   fullDocument,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}
}

// createErrorResponse crea una respuesta de error
func (c *InvoiceRequestConsumer) createErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	processingTime := time.Since(startTime).Milliseconds()

	return &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "softpymes",
		Status:         "error",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}
}
