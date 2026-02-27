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

const (
	queueName = "orders.events.inventory"
)

// ============================================
// STRUCTS LOCALES REPLICADOS
// (No importar del módulo orders — aislamiento)
// ============================================

// orderEventMessage replica la estructura de OrderEventMessage del módulo orders
type orderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *orderSnapshot         `json:"order"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// orderSnapshot replica OrderSnapshot del módulo orders
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

// orderItemSnapshot replica OrderItemSnapshot del módulo orders
type orderItemSnapshot struct {
	ProductID *string `json:"product_id,omitempty"`
	SKU       string  `json:"sku"`
	Quantity  int     `json:"quantity"`
}

// OrderConsumer consume eventos de órdenes desde RabbitMQ para mover inventario
type OrderConsumer struct {
	queue  rabbitmq.IQueue
	uc     app.IUseCase
	logger log.ILogger
}

// NewOrderConsumer crea un nuevo consumer de eventos de órdenes
func NewOrderConsumer(queue rabbitmq.IQueue, uc app.IUseCase, logger log.ILogger) *OrderConsumer {
	return &OrderConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("inventory.consumer"),
	}
}

// Start inicia el consumer en una goroutine
func (c *OrderConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, inventory order consumer disabled")
		return
	}

	// Declarar la cola (ya debería estar creada por orders/bundle.go, pero por seguridad)
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare inventory queue")
		return
	}

	c.logger.Info(ctx).Str("queue", queueName).Msg("Starting inventory order consumer")

	go func() {
		err := c.queue.Consume(ctx, queueName, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil // Siempre ACK — best-effort
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Inventory order consumer stopped with error")
		}
	}()
}

// handleMessage procesa un mensaje de evento de orden
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

	// Extraer items con ProductID
	items := c.extractItems(msg.Order.Items)
	if len(items) == 0 {
		return // Sin items con ProductID → nada que hacer
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

	default:
		// Eventos no relevantes para inventario
	}
}

func (c *OrderConsumer) handleReserve(ctx context.Context, msg orderEventMessage, businessID uint, items []dtos.OrderInventoryItem) {
	result, err := c.uc.ReserveStockForOrder(ctx, msg.OrderID, businessID, msg.Order.WarehouseID, items)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to reserve stock")
		return
	}
	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Bool("success", result.Success).
		Msg("Stock reserved for order")
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

// handleStatusChanged routea cambios de status a la acción correcta
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

// extractItems convierte items del snapshot a DTOs de inventario
// Solo incluye items que tienen ProductID (necesario para mover inventario)
func (c *OrderConsumer) extractItems(snapshotItems []orderItemSnapshot) []dtos.OrderInventoryItem {
	var items []dtos.OrderInventoryItem
	for _, si := range snapshotItems {
		if si.ProductID == nil || *si.ProductID == "" {
			continue
		}
		if si.Quantity <= 0 {
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
