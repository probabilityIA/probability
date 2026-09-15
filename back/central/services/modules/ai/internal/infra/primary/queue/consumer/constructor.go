package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Consumer struct {
	queue rabbitmq.IQueue
	uc    app.IUseCase
	log   log.ILogger
}

func New(queue rabbitmq.IQueue, uc app.IUseCase, logger log.ILogger) *Consumer {
	return &Consumer{queue: queue, uc: uc, log: logger}
}
