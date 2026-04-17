package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handlers struct {
	uc  app.IWebhookUseCase
	log log.ILogger
}

func New(uc app.IWebhookUseCase, logger log.ILogger) *Handlers {
	return &Handlers{
		uc:  uc,
		log: logger.WithModule("transport.envioclick.handler"),
	}
}
