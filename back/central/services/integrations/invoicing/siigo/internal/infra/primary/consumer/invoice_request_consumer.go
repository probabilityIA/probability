package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	siigoDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ═══════════════════════════════════════════════════════════════
// DTOs locales replicados del módulo Invoicing para deserialización
// (Regla de aislamiento: no importar entre módulos)
// ═══════════════════════════════════════════════════════════════

type invoiceCustomerData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	Address string `json:"address,omitempty"`
}

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
	Discount    float64  `json:"discount"`
}

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
	siigoClient       ports.ISiigoClient
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
}

// NewInvoiceRequestConsumer crea una nueva instancia del consumer
func NewInvoiceRequestConsumer(
	rabbit rabbitmq.IQueue,
	integrationCore integrationCore.IIntegrationService,
	siigoClient ports.ISiigoClient,
	responsePublisher *queue.ResponsePublisher,
	logger log.ILogger,
) *InvoiceRequestConsumer {
	return &InvoiceRequestConsumer{
		rabbit:            rabbit,
		integrationCore:   integrationCore,
		siigoClient:       siigoClient,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("siigo.invoice_request_consumer"),
	}
}

const (
	QueueSiigoRequests = "invoicing.siigo.requests"
)

// Start inicia el consumer
func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueSiigoRequests).
		Msg("Starting Siigo invoice request consumer")

	if err := c.rabbit.DeclareQueue(QueueSiigoRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueSiigoRequests, c.handleInvoiceRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	c.log.Info(ctx).
		Str("queue", QueueSiigoRequests).
		Msg("Siigo consumer started successfully")

	return nil
}

// handleInvoiceRequest procesa una solicitud de facturación
func (c *InvoiceRequestConsumer) handleInvoiceRequest(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

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
		Msg("Received Siigo invoice request")

	var response *queue.InvoiceResponseMessage
	switch request.Operation {
	case "create", "retry":
		response = c.processCreateInvoice(ctx, &request, startTime)
	default:
		c.log.Warn(ctx).
			Str("operation", request.Operation).
			Msg("Unknown operation")
		response = c.createErrorResponse(&request, "unknown_operation", "Unknown operation: "+request.Operation, startTime, nil)
	}

	if err := c.responsePublisher.PublishResponse(ctx, response); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Failed to publish response")
		return err
	}

	return nil
}

// processCreateInvoice procesa la creación de una factura en Siigo
func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	// 1. Obtener integration_id del DTO
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

	// 3. Desencriptar credenciales de Siigo
	username, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "username")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt username")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt username", startTime, nil)
	}

	accessKey, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "access_key")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt access_key")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt access_key", startTime, nil)
	}

	accountID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "account_id")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt account_id")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt account_id", startTime, nil)
	}

	partnerID, err := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "partner_id")
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to decrypt partner_id")
		return c.createErrorResponse(request, "decryption_failed", "Failed to decrypt partner_id", startTime, nil)
	}

	// api_url es opcional: si no está configurado, el cliente usa su default
	apiURL, _ := c.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_url")

	// 4. Combinar config de integración con config de facturación
	combinedConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		combinedConfig[k] = v
	}

	for k, v := range request.InvoiceData.Config {
		combinedConfig[k] = v
	}

	// 5. Construir request tipado para el cliente Siigo
	invoiceReq := &siigoDtos.CreateInvoiceRequest{
		Customer: siigoDtos.CustomerData{
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
		Credentials: siigoDtos.Credentials{
			Username:  username,
			AccessKey: accessKey,
			AccountID: accountID,
			PartnerID: partnerID,
			BaseURL:   apiURL,
		},
		Config: combinedConfig,
	}

	// 6. Llamar al cliente HTTP de Siigo
	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", request.InvoiceData.OrderID).
		Msg("Calling Siigo API")

	result, err := c.siigoClient.CreateInvoice(ctx, invoiceReq)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Siigo API call failed")

		var auditData *siigoDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createErrorResponse(request, "api_error", err.Error(), startTime, auditData)
	}

	// 7. Parsear issued_at si existe
	var issuedAt *time.Time
	if result.IssuedAt != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, result.IssuedAt); parseErr == nil {
			issuedAt = &parsed
		}
	}

	// 8. Construir response exitosa
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "siigo",
		Status:         "success",
		InvoiceNumber:  result.InvoiceNumber,
		ExternalID:     result.ExternalID,
		CUFE:           result.CUFE,
		IssuedAt:       issuedAt,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	if result.AuditData != nil {
		resp.AuditRequestURL = result.AuditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(result.AuditData.RequestPayload)
		resp.AuditResponseStatus = result.AuditData.ResponseStatus
		resp.AuditResponseBody = result.AuditData.ResponseBody
	}

	return resp
}

// mapItemsToClientDTOs convierte items del mensaje a DTOs del cliente Siigo
func mapItemsToClientDTOs(items []invoiceItemData) []siigoDtos.ItemData {
	result := make([]siigoDtos.ItemData, 0, len(items))
	for _, item := range items {
		result = append(result, siigoDtos.ItemData{
			ProductID:   item.ProductID,
			SKU:         item.SKU,
			Name:        item.Name,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			Tax:         item.Tax,
			TaxRate:     item.TaxRate,
			Discount:    item.Discount,
		})
	}
	return result
}

// createErrorResponse crea una respuesta de error
func (c *InvoiceRequestConsumer) createErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
	auditData *siigoDtos.AuditData,
) *queue.InvoiceResponseMessage {
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "siigo",
		Status:         "error",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
	}

	if auditData != nil {
		resp.AuditRequestURL = auditData.RequestURL
		resp.AuditRequestPayload = toMapPayload(auditData.RequestPayload)
		resp.AuditResponseStatus = auditData.ResponseStatus
		resp.AuditResponseBody = auditData.ResponseBody
	}

	return resp
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
