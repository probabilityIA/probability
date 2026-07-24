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

type invoiceCustomerData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	Address string `json:"address,omitempty"`
}

type invoiceItemData struct {
	ProductID       *string  `json:"product_id"`
	SKU             string   `json:"sku"`
	Name            string   `json:"name"`
	Description     *string  `json:"description"`
	Quantity        int      `json:"quantity"`
	UnitPrice       float64  `json:"unit_price"`
	UnitPriceBase   float64  `json:"unit_price_base"`
	TotalPrice      float64  `json:"total_price"`
	Tax             float64  `json:"tax"`
	TaxRate         *float64 `json:"tax_rate"`
	Discount        float64  `json:"discount"`
	DiscountPercent float64  `json:"discount_percent"`
	UnitPricePresentment     float64 `json:"unit_price_presentment"`
	UnitPriceBasePresentment float64 `json:"unit_price_base_presentment"`
	TotalPricePresentment    float64 `json:"total_price_presentment"`
	DiscountPresentment      float64 `json:"discount_presentment"`
	TaxPresentment           float64 `json:"tax_presentment"`
}

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

type InvoiceRequestMessage struct {
	InvoiceID     uint        `json:"invoice_id"`
	Provider      string      `json:"provider"`
	Operation     string      `json:"operation"`
	InvoiceData   invoiceData `json:"invoice_data"`
	CorrelationID string      `json:"correlation_id"`
	Timestamp     time.Time   `json:"timestamp"`
}

type InvoiceRequestConsumer struct {
	rabbit            rabbitmq.IQueue
	integrationCore   integrationCore.IIntegrationService
	softpymesClient   ports.ISoftpymesClient
	responsePublisher *queue.ResponsePublisher
	log               log.ILogger
	workers           int
}

func New(
	rabbit rabbitmq.IQueue,
	integrationCore integrationCore.IIntegrationService,
	softpymesClient ports.ISoftpymesClient,
	responsePublisher *queue.ResponsePublisher,
	logger log.ILogger,
	workers int,
) *InvoiceRequestConsumer {
	if workers < 1 {
		workers = 1
	}
	return &InvoiceRequestConsumer{
		rabbit:            rabbit,
		integrationCore:   integrationCore,
		softpymesClient:   softpymesClient,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("softpymes.invoice_request_consumer"),
		workers:           workers,
	}
}

const (
	QueueSoftpymesRequests = rabbitmq.QueueInvoicingSoftpymesRequests
)

func (c *InvoiceRequestConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		c.log.Warn(ctx).Msg("RabbitMQ client is nil, consumer cannot start")
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", QueueSoftpymesRequests).
		Int("workers", c.workers).
		Msg("Starting Softpymes invoice request consumer")

	if err := c.rabbit.DeclareQueue(QueueSoftpymesRequests, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return err
	}

	if err := c.rabbit.ConsumeConcurrent(ctx, QueueSoftpymesRequests, c.handleInvoiceRequest, c.workers); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return err
	}

	c.log.Info(ctx).
		Str("queue", QueueSoftpymesRequests).
		Int("workers", c.workers).
		Msg("Consumer started successfully")

	return nil
}

func (c *InvoiceRequestConsumer) handleInvoiceRequest(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var request InvoiceRequestMessage
	if err := json.Unmarshal(message, &request); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("Failed to unmarshal request - dropping corrupt message")
		return nil
	}

	c.log.Info(ctx).
		Uint("invoice_id", request.InvoiceID).
		Str("operation", request.Operation).
		Str("correlation_id", request.CorrelationID).
		Msg("Received invoice request")

	if request.Operation == "compare" {
		return c.processCompareRequest(ctx, &request)
	}

	if request.Operation == "list_items" {
		return c.processListItemsRequest(ctx, &request)
	}

	if request.Operation == "list_bank_accounts" {
		return c.processListBankAccountsRequest(ctx, &request)
	}

	if request.Operation == "reconcile_failed" {
		return c.processReconcileFailed(ctx, &request)
	}

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

	if err := c.responsePublisher.PublishResponse(ctx, response); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", request.InvoiceID).
			Msg("Failed to publish response")
		return err
	}

	return nil
}
