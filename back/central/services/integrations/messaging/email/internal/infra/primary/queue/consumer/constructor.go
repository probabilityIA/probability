package consumer

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type emailConsumer struct {
	rabbitMQ rabbitmq.IQueue
	useCase  ports.IEmailUseCase
	logger   log.ILogger
}

// New crea un nuevo consumer de emails desde RabbitMQ
func New(rabbitMQ rabbitmq.IQueue, useCase ports.IEmailUseCase, logger log.ILogger) *emailConsumer {
	return &emailConsumer{
		rabbitMQ: rabbitMQ,
		useCase:  useCase,
		logger:   logger,
	}
}
