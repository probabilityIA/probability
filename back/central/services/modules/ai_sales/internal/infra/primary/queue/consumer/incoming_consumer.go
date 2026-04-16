package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start declara la cola e inicia el consumo de mensajes
func (c *Consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueWhatsAppAIIncoming

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error declarando cola AI incoming")
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error iniciando consumer AI incoming")
		return fmt.Errorf("failed to start consumer %s: %w", queueName, err)
	}

	c.log.Info().Str("queue", queueName).Msg("AI Sales consumer iniciado")
	return nil
}

func (c *Consumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	var dto domain.IncomingMessageDTO
	if err := json.Unmarshal(messageBody, &dto); err != nil {
		c.log.Error().
			Err(err).
			Str("raw", string(messageBody)).
			Msg("Error deserializando mensaje AI incoming")
		return nil // No reintentar mensajes malformados
	}

	if dto.PhoneNumber == "" || dto.MessageText == "" {
		c.log.Warn().Msg("Mensaje AI incoming con campos vacios, descartando")
		return nil
	}

	c.log.Info().
		Str("phone", dto.PhoneNumber).
		Str("message_type", dto.MessageType).
		Msg("Procesando mensaje AI incoming")

	if err := c.useCase.HandleIncoming(ctx, dto); err != nil {
		c.log.Error().
			Err(err).
			Str("phone", dto.PhoneNumber).
			Msg("Error procesando mensaje AI incoming")
		// No retornar error para no reencolar - ya se envió mensaje de error al usuario
		return nil
	}

	return nil
}
