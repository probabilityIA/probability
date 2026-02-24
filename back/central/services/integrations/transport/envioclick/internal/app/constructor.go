package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCase handles EnvioClick transport operations
type UseCase struct {
	client domain.IEnvioClickClient
	log    log.ILogger
}

// New creates a new EnvioClick use case
func New(client domain.IEnvioClickClient, logger log.ILogger) *UseCase {
	return &UseCase{
		client: client,
		log:    logger.WithModule("transport.envioclick.usecase"),
	}
}
