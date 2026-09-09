package consumer

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IConsumer interface {
	Start(ctx context.Context) error
}

type consumer struct {
	queue   rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

func New(queue rabbitmq.IQueue, useCase app.IUseCase, logger log.ILogger) IConsumer {
	return &consumer{
		queue:   queue,
		useCase: useCase,
		log:     logger.WithModule("push-consumer"),
	}
}
