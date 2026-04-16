package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type orderEventMessage struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	OrderID   string `json:"order_id"`
}

type Consumer struct {
	queue   rabbitmq.IQueue
	logger  log.ILogger
	useCase ports.IScoreUseCase
}

func New(queue rabbitmq.IQueue, logger log.ILogger, useCase ports.IScoreUseCase) *Consumer {
	return &Consumer{queue: queue, logger: logger, useCase: useCase}
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := c.queue.DeclareQueue(rabbitmq.QueueOrdersToScore, true); err != nil {
		return err
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueOrdersToScore).Msg("Probability score consumer iniciado")
	return c.queue.Consume(ctx, rabbitmq.QueueOrdersToScore, c.handleMessage)
}

func (c *Consumer) handleMessage(msg []byte) error {
	var event orderEventMessage
	if err := json.Unmarshal(msg, &event); err != nil {
		c.logger.Error(context.Background()).Err(err).Msg("Error deserializando evento en probability consumer")
		return nil // ACK: mensaje malformado, no reintentar
	}

	// Solo procesar order.created y order.updated
	if event.EventType != "order.created" && event.EventType != "order.updated" {
		return nil // ACK: evento no relevante
	}

	if event.OrderID == "" {
		c.logger.Warn(context.Background()).Str("event_type", event.EventType).Msg("Evento sin order_id, ignorando")
		return nil
	}

	ctx := context.Background()
	if err := c.useCase.CalculateAndUpdateOrderScore(ctx, event.OrderID); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("order_id", event.OrderID).
			Str("event_type", event.EventType).
			Msg("Error calculando score de la orden")
		return nil // ACK: no reintentar para evitar loops, el error ya se logueó
	}

	return nil
}
