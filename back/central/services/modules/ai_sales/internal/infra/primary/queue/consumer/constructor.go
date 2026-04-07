package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	queue   rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

func New(queue rabbitmq.IQueue, useCase app.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{
		queue:   queue,
		useCase: useCase,
		log:     logger,
	}
}
