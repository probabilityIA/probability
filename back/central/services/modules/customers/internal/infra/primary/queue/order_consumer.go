package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const queueName = rabbitmq.QueueOrdersToCustomers

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

type orderSnapshot struct {
	ID                  string              `json:"id"`
	OrderNumber         string              `json:"order_number"`
	TotalAmount         float64             `json:"total_amount"`
	Currency            string              `json:"currency"`
	Platform            string              `json:"platform"`
	CustomerID          *uint               `json:"customer_id,omitempty"`
	CustomerName        string              `json:"customer_name"`
	CustomerEmail       string              `json:"customer_email,omitempty"`
	CustomerPhone       string              `json:"customer_phone,omitempty"`
	CustomerDNI         string              `json:"customer_dni,omitempty"`
	ShippingStreet      string              `json:"shipping_street,omitempty"`
	ShippingCity        string              `json:"shipping_city,omitempty"`
	ShippingState       string              `json:"shipping_state,omitempty"`
	ShippingCountry     string              `json:"shipping_country,omitempty"`
	ShippingPostalCode  string              `json:"shipping_postal_code,omitempty"`
	IsPaid              bool                `json:"is_paid"`
	DeliveryProbability float64             `json:"delivery_probability,omitempty"`
	Status              string              `json:"status,omitempty"`
	Items               []orderItemSnapshot `json:"items,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
}

type orderItemSnapshot struct {
	ProductID *string `json:"product_id,omitempty"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	ImageURL  *string `json:"image_url,omitempty"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
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
		logger: logger.WithModule("customers.consumer"),
	}
}

func (c *OrderConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, customers order consumer disabled")
		return
	}

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare customers queue")
		return
	}

	c.logger.Info(ctx).Str("queue", queueName).Msg("Starting customers order consumer")

	go func() {
		err := c.queue.Consume(ctx, queueName, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Customers order consumer stopped with error")
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

	businessID := uint(0)
	if msg.BusinessID != nil {
		businessID = *msg.BusinessID
	}
	if businessID == 0 {
		c.logger.Warn(ctx).Str("order_id", msg.OrderID).Msg("Order event without business_id, skipping")
		return
	}

	previousStatus, _ := msg.Changes["previous_status"].(string)
	currentStatus, _ := msg.Changes["current_status"].(string)

	items := make([]dtos.OrderEventItemDTO, 0, len(msg.Order.Items))
	for _, item := range msg.Order.Items {
		items = append(items, dtos.OrderEventItemDTO{
			ProductID:    item.ProductID,
			ProductName:  item.Name,
			ProductSKU:   item.SKU,
			ProductImage: item.ImageURL,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
		})
	}

	event := dtos.OrderEventDTO{
		EventType:           msg.EventType,
		OrderID:             msg.OrderID,
		BusinessID:          businessID,
		CustomerID:          msg.Order.CustomerID,
		CustomerName:        msg.Order.CustomerName,
		CustomerEmail:       msg.Order.CustomerEmail,
		CustomerPhone:       msg.Order.CustomerPhone,
		CustomerDNI:         msg.Order.CustomerDNI,
		TotalAmount:         msg.Order.TotalAmount,
		Currency:            msg.Order.Currency,
		Platform:            msg.Order.Platform,
		Status:              msg.Order.Status,
		IsPaid:              msg.Order.IsPaid,
		DeliveryProbability: msg.Order.DeliveryProbability,
		ShippingStreet:      msg.Order.ShippingStreet,
		ShippingCity:        msg.Order.ShippingCity,
		ShippingState:       msg.Order.ShippingState,
		ShippingCountry:     msg.Order.ShippingCountry,
		ShippingPostalCode:  msg.Order.ShippingPostalCode,
		OrderNumber:         msg.Order.OrderNumber,
		OrderedAt:           msg.Order.CreatedAt,
		Items:               items,
		PreviousStatus:      previousStatus,
		CurrentStatus:       currentStatus,
	}

	c.logger.Info(ctx).
		Str("event_type", msg.EventType).
		Str("order_id", msg.OrderID).
		Uint("business_id", businessID).
		Msg("Processing order event for customers")

	if err := c.uc.ProcessOrderEvent(ctx, event); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("order_id", msg.OrderID).
			Msg("Failed to process order event for customers")
	}
}
