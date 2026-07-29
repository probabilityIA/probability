package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Handlers struct {
	uc     app.IWebhookUseCase
	log    log.ILogger
	rabbit rabbitmq.IQueue
}

func New(uc app.IWebhookUseCase, logger log.ILogger, rabbit rabbitmq.IQueue) *Handlers {
	return &Handlers{
		uc:     uc,
		log:    logger.WithModule("transport.shipit.handler"),
		rabbit: rabbit,
	}
}
