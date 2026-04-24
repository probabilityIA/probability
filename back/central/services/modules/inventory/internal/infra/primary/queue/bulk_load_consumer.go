package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// BulkLoadConsumer consume mensajes de carga masiva de inventario desde RabbitMQ
type BulkLoadConsumer struct {
	queue          rabbitmq.IQueue
	uc             app.IUseCase
	eventPublisher ports.IInventoryEventPublisher
	logger         log.ILogger
}

// NewBulkLoadConsumer crea un nuevo consumer de carga masiva
func NewBulkLoadConsumer(queue rabbitmq.IQueue, uc app.IUseCase, eventPublisher ports.IInventoryEventPublisher, logger log.ILogger) *BulkLoadConsumer {
	return &BulkLoadConsumer{
		queue:          queue,
		uc:             uc,
		eventPublisher: eventPublisher,
		logger:         logger.WithModule("inventory.bulk_load_consumer"),
	}
}

// Start inicia el consumer en una goroutine
func (c *BulkLoadConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, bulk load consumer disabled")
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueInventoryBulkLoad, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare bulk load queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueInventoryBulkLoad).Msg("Starting bulk load consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueInventoryBulkLoad, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil // Siempre ACK
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Bulk load consumer stopped with error")
		}
	}()
}

func (c *BulkLoadConsumer) handleMessage(ctx context.Context, body []byte) {
	var dto request.BulkLoadDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal bulk load message")
		return
	}

	c.logger.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Uint("warehouse_id", dto.WarehouseID).
		Int("items", len(dto.Items)).
		Msg("Processing bulk load request")

	result, err := c.uc.BulkLoadInventory(ctx, dto)
	if err != nil {
		c.logger.Error(ctx).Err(err).
			Uint("business_id", dto.BusinessID).
			Msg("Bulk load failed")

		// Publicar evento de error via SSE
		if c.eventPublisher != nil {
			_ = c.eventPublisher.PublishInventoryEvent(ctx, ports.InventoryEvent{
				EventType:   "bulk_load.failed",
				BusinessID:  dto.BusinessID,
				WarehouseID: dto.WarehouseID,
				Data: map[string]any{
					"error":       err.Error(),
					"total_items": len(dto.Items),
				},
			})
		}
		return
	}

	c.logger.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Int("total", result.TotalItems).
		Int("success", result.SuccessCount).
		Int("failures", result.FailureCount).
		Msg("Bulk load completed")

	// Publicar evento de resultado via SSE
	if c.eventPublisher != nil {
		_ = c.eventPublisher.PublishInventoryEvent(ctx, ports.InventoryEvent{
			EventType:   "bulk_load.completed",
			BusinessID:  dto.BusinessID,
			WarehouseID: dto.WarehouseID,
			Data: map[string]any{
				"total_items":   result.TotalItems,
				"success_count": result.SuccessCount,
				"failure_count": result.FailureCount,
			},
		})
	}
}
