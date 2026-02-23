package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	factDtos "github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ═══════════════════════════════════════════════════════════════
// DTOs locales — structs de deserialización del mensaje RabbitMQ
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

// InvoiceRequestMessage es el mensaje recibido desde el Invoicing Module
type InvoiceRequestMessage struct {
	InvoiceID     uint        `json:"invoice_id"`
	Provider      string      `json:"provider"`
	Operation     string      `json:"operation"`
	InvoiceData   invoiceData `json:"invoice_data"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"timestamp"`
}

// InvoiceRequestConsumer consume solicitudes de facturación desde RabbitMQ
// y delega toda la lógica de negocio al use case.
type InvoiceRequestConsumer struct {
	rabbit            rabbitmq.IQueue
	useCase           ports.IInvoiceUseCase
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
}

// NewInvoiceRequestConsumer crea una nueva instancia del consumer.
// Solo recibe el use case — no adapters secundarios directamente.
func NewInvoiceRequestConsumer(
	rabbit rabbitmq.IQueue,
	useCase ports.IInvoiceUseCase,
	responsePublisher *queue.ResponsePublisher,
	logger log.ILogger,
) *InvoiceRequestConsumer {
	return &InvoiceRequestConsumer{
		rabbit:            rabbit,
		useCase:           useCase,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("factus.invoice_request_consumer"),
	}
}

const (
	QueueFactusRequests = "invoicing.factus.requests"
)

// Start inicia el consumer
func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueFactusRequests).
		Msg("Starting Factus invoice request consumer")

	if err := c.rabbit.DeclareQueue(QueueFactusRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueFactusRequests, c.handleInvoiceRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	c.log.Info(ctx).
		Str("queue", QueueFactusRequests).
		Msg("Factus consumer started successfully")

	return nil
}

// handleInvoiceRequest deserializa el mensaje y despacha al handler correcto
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
		Msg("Received Factus invoice request")

	var response *queue.InvoiceResponseMessage
	switch request.Operation {
	case "create", "retry":
		response = c.processCreateInvoice(ctx, &request, startTime)
	default:
		c.log.Warn(ctx).Str("operation", request.Operation).Msg("Unknown operation")
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

// processCreateInvoice construye el DTO de dominio y delega al use case.
// No contiene lógica de negocio — solo traducción del mensaje al dominio.
func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	if request.InvoiceData.IntegrationID == 0 {
		c.log.Error(ctx).Msg("integration_id is 0 in invoice_data")
		return c.createErrorResponse(request, "missing_integration_id", "integration_id is 0", startTime, nil)
	}

	// Construir el DTO de dominio para el use case
	req := &factDtos.ProcessInvoiceRequest{
		InvoiceID:     request.InvoiceID,
		Operation:     request.Operation,
		CorrelationID: request.CorrelationID,
		IntegrationID: request.InvoiceData.IntegrationID,
		Customer: factDtos.CustomerData{
			Name:    request.InvoiceData.Customer.Name,
			Email:   request.InvoiceData.Customer.Email,
			Phone:   request.InvoiceData.Customer.Phone,
			DNI:     request.InvoiceData.Customer.DNI,
			Address: request.InvoiceData.Customer.Address,
		},
		Items:        mapItemsToDomain(request.InvoiceData.Items),
		Total:        request.InvoiceData.Total,
		Subtotal:     request.InvoiceData.Subtotal,
		Tax:          request.InvoiceData.Tax,
		Discount:     request.InvoiceData.Discount,
		ShippingCost: request.InvoiceData.ShippingCost,
		Currency:     request.InvoiceData.Currency,
		OrderID:      request.InvoiceData.OrderID,
		Config:       request.InvoiceData.Config,
	}

	// Delegar toda la lógica de negocio al use case
	result, err := c.useCase.CreateInvoice(ctx, req)
	if err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Use case returned error")

		var auditData *factDtos.AuditData
		if result != nil {
			auditData = result.AuditData
		}
		return c.createErrorResponse(request, "processing_error", err.Error(), startTime, auditData)
	}

	// Parsear issued_at
	var issuedAt *time.Time
	if result.IssuedAt != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, result.IssuedAt); parseErr == nil {
			issuedAt = &parsed
		}
	}

	// Construir respuesta exitosa
	processingTime := time.Since(startTime).Milliseconds()
	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "factus",
		Status:         "success",
		InvoiceNumber:  result.InvoiceNumber,
		ExternalID:     result.ExternalID,
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

// mapItemsToDomain convierte items del mensaje RabbitMQ a DTOs de dominio
func mapItemsToDomain(items []invoiceItemData) []factDtos.ItemData {
	result := make([]factDtos.ItemData, 0, len(items))
	for _, item := range items {
		result = append(result, factDtos.ItemData{
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

// createErrorResponse construye una respuesta de error para la queue de respuestas
func (c *InvoiceRequestConsumer) createErrorResponse(
	request *InvoiceRequestMessage,
	errorCode string,
	errorMsg string,
	startTime time.Time,
	auditData *factDtos.AuditData,
) *queue.InvoiceResponseMessage {
	processingTime := time.Since(startTime).Milliseconds()

	resp := &queue.InvoiceResponseMessage{
		InvoiceID:      request.InvoiceID,
		Provider:       "factus",
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

// toMapPayload convierte cualquier valor a map[string]interface{} via JSON
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
