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
	UnitPriceBase float64 `json:"unit_price_base"`
	TotalPrice  float64  `json:"total_price"`
	Tax         float64  `json:"tax"`
	TaxRate     *float64 `json:"tax_rate"`
	Discount        float64  `json:"discount"`
	DiscountPercent float64  `json:"discount_percent"`
	// Precios en moneda presentment (moneda local, ej: COP)
	UnitPricePresentment      float64 `json:"unit_price_presentment"`
	UnitPriceBasePresentment  float64 `json:"unit_price_base_presentment"`
	TotalPricePresentment     float64 `json:"total_price_presentment"`
	DiscountPresentment       float64 `json:"discount_presentment"`
	TaxPresentment            float64 `json:"tax_presentment"`
}

// invoiceData datos completos (replicado de invoicing module)
type invoiceData struct {
	IntegrationID    uint                   `json:"integration_id"`
	Customer         invoiceCustomerData    `json:"customer"`
	Items            []invoiceItemData      `json:"items"`
	Total            float64                `json:"total"`
	Subtotal         float64                `json:"subtotal"`
	Tax              float64                `json:"tax"`
	Discount         float64                `json:"discount"`
	ShippingCost     float64                `json:"shipping_cost"`
	ShippingDiscount float64                `json:"shipping_discount"`
	ShippingCostBase float64                `json:"shipping_cost_base"`
	Currency         string                 `json:"currency"`
	OrderID          string                 `json:"order_id"`
	OrderNumber      string                 `json:"order_number,omitempty"`
	Config           map[string]interface{} `json:"config"`
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

	// Operación "list_items": flujo propio (publica ListItemsResponseMessage, no InvoiceResponseMessage)
	if request.Operation == "list_items" {
		return c.processListItemsRequest(ctx, &request)
	}

	// Operación "list_bank_accounts": flujo propio (publica ListBankAccountsResponseMessage)
	if request.Operation == "list_bank_accounts" {
		return c.processListBankAccountsRequest(ctx, &request)
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
	case "cash_receipt":
		response = c.processCashReceipt(ctx, &request, startTime)
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
