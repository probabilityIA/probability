package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	repository ports.IRepository
	logger     log.ILogger
}

// New crea una nueva instancia del caso de uso
func New(repository ports.IRepository, logger log.ILogger) ports.IUseCase {
	return &useCase{
		repository: repository,
		logger:     logger.WithModule("notification_config_usecase"),
	}
}
