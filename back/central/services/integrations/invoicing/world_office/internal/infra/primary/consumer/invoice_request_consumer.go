package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/infra/secondary/queue"
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
	worldOfficeClient ports.IWorldOfficeClient
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
}

// NewInvoiceRequestConsumer crea una nueva instancia del consumer
func NewInvoiceRequestConsumer(
	rabbit rabbitmq.IQueue,
	integrationCore integrationCore.IIntegrationService,
	worldOfficeClient ports.IWorldOfficeClient,
	responsePublisher *queue.ResponsePublisher,
	logger log.ILogger,
) *InvoiceRequestConsumer {
	return &InvoiceRequestConsumer{
		rabbit:            rabbit,
		integrationCore:   integrationCore,
		worldOfficeClient: worldOfficeClient,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("world_office.invoice_request_consumer"),
	}
}

const (
	QueueWorldOfficeRequests = "invoicing.world_office.requests"
)

// Start inicia el consumer
func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueWorldOfficeRequests).
		Msg("Starting World Office invoice request consumer")

	if err := c.rabbit.DeclareQueue(QueueWorldOfficeRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.Consume(ctx, QueueWorldOfficeRequests, c.handleInvoiceRequest); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	c.log.Info(ctx).
		Str("queue", QueueWorldOfficeRequests).
		Msg("World Office consumer started successfully")

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
		Msg("Received World Office invoice request")

	var response *queue.InvoiceResponseMessage
	switch request.Operation {
	case "create", "retry":
		response = c.processCreateInvoice(ctx, &request, startTime)
	default:
		c.log.Warn(ctx).
			Str("operation", request.Operation).
			Msg("Unknown operation")
		response = c.createErrorResponse(&request, "unknown_operation", "Unknown operation: "+request.Operation, startTime)
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

// processCreateInvoice procesa la creación de una factura en World Office
// TODO: Implementar lógica real de creación de factura con la API de World Office
func (c *InvoiceRequestConsumer) processCreateInvoice(
	ctx context.Context,
	request *InvoiceRequestMessage,
	startTime time.Time,
) *queue.InvoiceResponseMessage {
	c.log.Warn(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("order_id", request.InvoiceData.OrderID).
		Msg("⚠️ World Office invoice creation not yet implemented")
	return c.createErrorResponse(request, "not_implemented", "world_office: invoice creation not yet implemented", startTime)
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
		Provider:       "world_office",
		Status:         "error",
		Error:          errorMsg,
		ErrorCode:      errorCode,
		CorrelationID:  request.CorrelationID,
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
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
