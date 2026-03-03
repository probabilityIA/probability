package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	exchangeName = rabbitmq.ExchangeInventory
	exchangeType = "topic"
)

// SyncPublisher publica mensajes de sync de inventario a RabbitMQ
type SyncPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

// New crea un nuevo publisher y declara el exchange
func New(queue rabbitmq.IQueue, logger log.ILogger) ports.ISyncPublisher {
	if queue == nil {
		return &SyncPublisher{queue: nil, logger: logger}
	}

	// Declarar exchange
	if err := queue.DeclareExchange(exchangeName, exchangeType, true); err != nil {
		logger.Error().Err(err).Msg("Failed to declare inventory exchange")
	}

	return &SyncPublisher{queue: queue, logger: logger}
}

func (p *SyncPublisher) PublishInventorySync(ctx context.Context, msg ports.InventorySyncMessage) error {
	if p.queue == nil {
		return nil
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory sync message: %w", err)
	}

	routingKey := fmt.Sprintf("sync.%d", msg.IntegrationID)

	if err := p.queue.PublishToExchange(ctx, exchangeName, routingKey, body); err != nil {
		p.logger.Error().
			Err(err).
			Str("product_id", msg.ProductID).
			Uint("integration_id", msg.IntegrationID).
			Msg("Failed to publish inventory sync message")
		return err
	}

	p.logger.Info().
		Str("product_id", msg.ProductID).
		Uint("integration_id", msg.IntegrationID).
		Int("new_quantity", msg.NewQuantity).
		Str("source", msg.Source).
		Msg("Inventory sync message published")

	return nil
}
