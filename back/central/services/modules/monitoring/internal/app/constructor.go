package app

import (
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	publisher ports.IAlertPublisher
	log       log.ILogger
	env       env.IConfig
}

// New crea una nueva instancia del use case de monitoreo
func New(publisher ports.IAlertPublisher, logger log.ILogger, environment env.IConfig) ports.IUseCase {
	return &useCase{
		publisher: publisher,
		log:       logger,
		env:       environment,
	}
}
