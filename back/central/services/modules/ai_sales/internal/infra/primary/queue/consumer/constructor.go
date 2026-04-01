package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Consumer consume mensajes de la cola whatsapp.ai.incoming
type Consumer struct {
	queue   rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer de mensajes AI incoming
func New(queue rabbitmq.IQueue, useCase app.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		queue:   queue,
		useCase: useCase,
		log:     logger,
	}
}
