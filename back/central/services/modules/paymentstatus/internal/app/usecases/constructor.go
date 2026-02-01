package usecases

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea una nueva instancia del caso de uso
func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{
		repo:   repo,
		logger: logger.WithModule("usecases"),
	}
}
