package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type inventoryFeedbackMessage struct {
	OrderID    string `json:"order_id"`
	BusinessID uint   `json:"business_id"`
	Success    bool   `json:"success"`
	EventType  string `json:"event_type"`
}

type InventoryConsumer struct {
	queue           rabbitmq.IQueue
	repo            ports.IRepository
	rabbitPublisher ports.IOrderRabbitPublisher
	logger          log.ILogger
}

func NewInventoryConsumer(queue rabbitmq.IQueue, repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher, logger log.ILogger) *InventoryConsumer {
	return &InventoryConsumer{
		queue:           queue,
		repo:            repo,
		rabbitPublisher: rabbitPublisher,
		logger:          logger.WithModule("orders.inventory.consumer"),
	}
}

func (c *InventoryConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, inventory feedback consumer disabled")
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueInventoryOrderFeedback, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare inventory feedback queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueInventoryOrderFeedback).Msg("Starting inventory feedback consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueInventoryOrderFeedback, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Inventory feedback consumer stopped with error")
		}
	}()
}

func (c *InventoryConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg inventoryFeedbackMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal inventory feedback message")
		return
	}

	if msg.OrderID == "" {
		return
	}

	targetCode := "picking"
	if !msg.Success {
		targetCode = "inventory_issue"
	}

	c.logger.Info(ctx).
		Str("order_id", msg.OrderID).
		Uint("business_id", msg.BusinessID).
		Bool("success", msg.Success).
		Str("target_status", targetCode).
		Msg("Inventory feedback received - updating order status")

	statusID, err := c.repo.GetOrderStatusIDByCode(ctx, targetCode)
	if err != nil || statusID == nil {
		c.logger.Warn(ctx).Str("order_id", msg.OrderID).Str("status", targetCode).Msg("status not found, skipping")
		return
	}

	order, err := c.repo.GetOrderByID(ctx, msg.OrderID)
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Msg("Failed to get order before status update")
		return
	}

	previousStatus := order.Status
	if previousStatus == targetCode {
		return
	}

	if err := c.repo.UpdateOrderStatus(ctx, msg.OrderID, targetCode, statusID); err != nil {
		c.logger.Error(ctx).Err(err).Str("order_id", msg.OrderID).Str("status", targetCode).Msg("Failed to update order status")
		return
	}

	if c.rabbitPublisher == nil {
		return
	}

	order.Status = targetCode
	order.StatusID = statusID
	if err := c.rabbitPublisher.PublishOrderStatusChanged(ctx, order, previousStatus, targetCode); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("order_id", msg.OrderID).
			Str("previous_status", previousStatus).
			Str("current_status", targetCode).
			Msg("Failed to publish order.status_changed event after inventory feedback")
	}
}
