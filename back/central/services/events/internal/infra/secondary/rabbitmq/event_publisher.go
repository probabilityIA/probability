package rabbitmq

import (
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Re-export constantes desde shared para uso interno del módulo
const (
	ExchangeName = rabbitmq.EventsExchangeName
	QueueName    = rabbitmq.EventsQueueName
)

// SetupInfrastructure declara exchange, queue y binding para el sistema de eventos.
// Se llama una vez al inicializar el módulo.
func SetupInfrastructure(rabbitMQ rabbitmq.IQueue, logger log.ILogger) {
	// Declarar exchange tipo topic (durable)
	if err := rabbitMQ.DeclareExchange(ExchangeName, "topic", true); err != nil {
		logger.Error().
			Err(err).
			Str("exchange", ExchangeName).
			Msg("Error declarando events exchange")
	}

	// Declarar queue durable
	if err := rabbitMQ.DeclareQueue(QueueName, true); err != nil {
		logger.Error().
			Err(err).
			Str("queue", QueueName).
			Msg("Error declarando events queue")
	}

	// Bind queue con wildcard "#" para recibir todos los routing keys
	if err := rabbitMQ.BindQueue(QueueName, ExchangeName, "#"); err != nil {
		logger.Error().
			Err(err).
			Str("queue", QueueName).
			Str("exchange", ExchangeName).
			Msg("Error bindeando events queue a exchange")
	}

	logger.Info().
		Str("exchange", ExchangeName).
		Str("queue", QueueName).
		Msg("Events RabbitMQ infrastructure declarada")
}
