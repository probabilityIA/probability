package rabbitmq

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Re-export constantes desde shared para uso interno del módulo
const (
	ExchangeName = rabbitmq.EventsExchangeName
	QueueName    = rabbitmq.EventsQueueName
)

// SetupInfrastructure declara exchange, queue y binding para el sistema de eventos.
// Se llama una vez al inicializar el módulo. Retorna error si falla la declaración.
func SetupInfrastructure(rabbitMQ rabbitmq.IQueue, logger log.ILogger) error {
	ctx := context.Background()

	// Declarar exchange tipo topic (durable)
	if err := rabbitMQ.DeclareExchange(ExchangeName, "topic", true); err != nil {
		logger.Error(ctx).
			Err(err).
			Str("exchange", ExchangeName).
			Msg("Error declarando events exchange")
		return fmt.Errorf("error declarando events exchange: %w", err)
	}

	// Declarar queue durable
	if err := rabbitMQ.DeclareQueue(QueueName, true); err != nil {
		logger.Error(ctx).
			Err(err).
			Str("queue", QueueName).
			Msg("Error declarando events queue")
		return fmt.Errorf("error declarando events queue: %w", err)
	}

	// Bind queue con wildcard "#" para recibir todos los routing keys
	if err := rabbitMQ.BindQueue(QueueName, ExchangeName, "#"); err != nil {
		logger.Error(ctx).
			Err(err).
			Str("queue", QueueName).
			Str("exchange", ExchangeName).
			Msg("Error bindeando events queue a exchange")
		return fmt.Errorf("error bindeando events queue a exchange: %w", err)
	}

	logger.Info(ctx).
		Str("exchange", ExchangeName).
		Str("queue", QueueName).
		Msg("Events RabbitMQ infrastructure declarada")

	return nil
}
