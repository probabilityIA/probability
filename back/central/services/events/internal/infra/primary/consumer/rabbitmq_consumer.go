package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	rmqPublisher "github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// RabbitMQConsumer consume eventos de la queue unificada y los despacha
type RabbitMQConsumer struct {
	rabbitMQ   rabbitmq.IQueue
	dispatcher ports.IEventDispatcher
	logger     log.ILogger
}

// New crea un nuevo RabbitMQConsumer
func New(
	rabbitMQ rabbitmq.IQueue,
	dispatcher ports.IEventDispatcher,
	logger log.ILogger,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		rabbitMQ:   rabbitMQ,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

// Start inicia el consumer en background
func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rmqPublisher.QueueName).
		Msg("Iniciando consumer de eventos unificado")

	return c.rabbitMQ.Consume(ctx, rmqPublisher.QueueName, func(body []byte) error {
		return c.handleMessage(ctx, body)
	})
}

// handleMessage deserializa un mensaje y lo despacha
func (c *RabbitMQConsumer) handleMessage(ctx context.Context, body []byte) error {
	var envelope dtos.EventEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("body", string(body)).
			Msg("Error deserializando evento de RabbitMQ")
		// No requeue mensajes malformados
		return nil
	}

	event := entities.Event{
		ID:            envelope.ID,
		Type:          envelope.Type,
		Category:      envelope.Category,
		BusinessID:    envelope.BusinessID,
		IntegrationID: envelope.IntegrationID,
		Timestamp:     envelope.Timestamp,
		Data:          envelope.Data,
		Metadata:      envelope.Metadata,
	}

	c.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("business_id", event.BusinessID).
		Msg("Evento recibido de RabbitMQ, despachando")

	return c.dispatcher.HandleEvent(ctx, event)
}
