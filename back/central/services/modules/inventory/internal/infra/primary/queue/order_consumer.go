package queue

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const queueName = rabbitmq.QueueOrdersToInventory

type orderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *orderSnapshot         `json:"order"`
	Changes       map[string]any `json:"changes,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

type orderSnapshot struct {
	ID                  string              `json:"id"`
	OrderNumber         string              `json:"order_number"`
	TotalAmount         float64             `json:"total_amount"`
	IntegrationID       uint                `json:"integration_id"`
	Platform            string              `json:"platform"`
	OrderStatusID       *uint               `json:"order_status_id,omitempty"`
	FulfillmentStatusID *uint               `json:"fulfillment_status_id,omitempty"`
	WarehouseID         *uint               `json:"warehouse_id,omitempty"`
	Items               []orderItemSnapshot `json:"items,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
}

type orderItemSnapshot struct {
	ProductID *string `json:"product_id,omitempty"`
	SKU       string  `json:"sku"`
	Quantity  int     `json:"quantity"`
}

type inventoryFeedbackMessage struct {
	OrderID    string `json:"order_id"`
	BusinessID uint   `json:"business_id"`
	Success    bool   `json:"success"`
	EventType  string `json:"event_type"`
}

type OrderConsumer struct {
	queue  rabbitmq.IQueue
	uc     app.IUseCase
	logger log.ILogger
}

func NewOrderConsumer(queue rabbitmq.IQueue, uc app.IUseCase, logger log.ILogger) *OrderConsumer {
	return &OrderConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("inventory.consumer"),
	}
}

func (c *OrderConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, inventory order consumer disabled")
		return
	}

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare inventory queue")
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueInventoryOrderFeedback, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare inventory feedback queue")
	}

	c.logger.Info(ctx).Str("queue", queueName).Msg("Starting inventory order consumer")

	go func() {
		err := c.queue.Consume(ctx, queueName, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Inventory order consumer stopped with error")
		}
	}()
}

func (c *OrderConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg orderEventMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal order event message")
		return
	}

	if msg.Order == nil {
		c.logger.Warn(ctx).Str("event_type", msg.EventType).Msg("Order event without snapshot, skipping")
		return
	}

	items := c.extractItems(msg.Order.Items)
	if len(items) == 0 {
		c.logger.Warn(ctx).
			Str("event_type", msg.EventType).
			Str("order_id", msg.OrderID).
			Int("raw_items", len(msg.Order.Items)).
			Msg("No trackable items in order event, skipping inventory")
		return
	}

	businessID := uint(0)
	if msg.BusinessID != nil {
		businessID = *msg.BusinessID
	}
	if businessID == 0 {
		c.logger.Warn(ctx).Str("order_id", msg.OrderID).Msg("Order event without business_id, skipping")
		return
	}

	c.logger.Info(ctx).
		Str("event_type", msg.EventType).
		Str("order_id", msg.OrderID).
		Uint("business_id", businessID).
		Int("items", len(items)).
		Msg("Processing order event for inventory")

	switch msg.EventType {
	case "order.created":
		c.handleReserve(ctx, msg, businessID, items)

	case "order.cancelled":
		c.handleRelease(ctx, msg, businessID, items)

	case "order.shipped", "order.completed":
		c.handleConfirmSale(ctx, msg, businessID, items)

	case "order.refunded":
		c.handleReturn(ctx, msg, businessID, items)

	case "order.status_changed":
		c.handleStatusChanged(ctx, msg, businessID, items)
	}
}

func (c *OrderConsumer) handleReserve(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	result, err := c.uc.ReserveStockForOrder(ctx, msg.OrderID, businessID, msg.Order.WarehouseID, items)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to reserve stock")
		c.publishFeedback(msg.OrderID, businessID, false)
		return
	}

	allSufficient := true
	for _, item := range result.ItemResults {
		if !item.Sufficient {
			allSufficient = false
			break
		}
	}

	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Bool("all_sufficient", allSufficient).
		Msg("Stock reserve result for order")

	c.publishFeedback(msg.OrderID, businessID, allSufficient)
}

func (c *OrderConsumer) handleRelease(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	result, err := c.uc.ReleaseStockForOrder(ctx, msg.OrderID, businessID, msg.Order.WarehouseID, items)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to release stock")
		return
	}
	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Bool("success", result.Success).
		Msg("Stock released for cancelled order")
}

func (c *OrderConsumer) handleConfirmSale(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	result, err := c.uc.ConfirmSaleForOrder(ctx, msg.OrderID, businessID, msg.Order.WarehouseID, items)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to confirm sale")
		return
	}
	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Bool("success", result.Success).
		Msg("Sale confirmed for order")
}

func (c *OrderConsumer) handleReturn(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	result, err := c.uc.ReturnStockForOrder(ctx, msg.OrderID, businessID, msg.Order.WarehouseID, items)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to return stock")
		return
	}
	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Bool("success", result.Success).
		Msg("Stock returned for refunded order")
}

func (c *OrderConsumer) handleStatusChanged(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	currentStatus, _ := msg.Changes["current_status"].(string)
	if currentStatus == "" {
		return
	}

	lower := strings.ToLower(currentStatus)

	switch {
	case strings.Contains(lower, "shipped") || strings.Contains(lower, "completed") || strings.Contains(lower, "delivered"):
		c.handleConfirmSale(ctx, msg, businessID, items)
	case strings.Contains(lower, "cancelled") || strings.Contains(lower, "canceled"):
		c.handleRelease(ctx, msg, businessID, items)
	case strings.Contains(lower, "refund"):
		c.handleReturn(ctx, msg, businessID, items)
	}
}

func (c *OrderConsumer) publishFeedback(orderID string, businessID uint, success bool) {
	if c.queue == nil {
		return
	}

	eventType := "inventory.reserved"
	if !success {
		eventType = "inventory.insufficient"
	}

	msg := inventoryFeedbackMessage{
		OrderID:    orderID,
		BusinessID: businessID,
		Success:    success,
		EventType:  eventType,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	go func() {
		_ = c.queue.Publish(context.Background(), rabbitmq.QueueInventoryOrderFeedback, body)
	}()
}

func (c *OrderConsumer) extractItems(snapshotItems []orderItemSnapshot) []dtos.OrderInventoryItem {
	var items []dtos.OrderInventoryItem
	for i, si := range snapshotItems {
		productIDStr := ""
		if si.ProductID != nil {
			productIDStr = *si.ProductID
		}
		if si.ProductID == nil || *si.ProductID == "" {
			c.logger.Warn(context.Background()).
				Int("item_index", i).
				Str("sku", si.SKU).
				Int("quantity", si.Quantity).
				Str("product_id", productIDStr).
				Msg("Item skipped: product_id is nil or empty")
			continue
		}
		if si.Quantity <= 0 {
			c.logger.Warn(context.Background()).
				Int("item_index", i).
				Str("sku", si.SKU).
				Str("product_id", productIDStr).
				Int("quantity", si.Quantity).
				Msg("Item skipped: quantity is zero or negative")
			continue
		}
		items = append(items, dtos.OrderInventoryItem{
			ProductID: *si.ProductID,
			SKU:       si.SKU,
			Quantity:  si.Quantity,
		})
	}
	return items
}
