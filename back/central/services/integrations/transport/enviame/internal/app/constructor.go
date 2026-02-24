package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCase implements the transport operations for Enviame
type UseCase struct {
	client domain.IEnviameClient
	log    log.ILogger
}

// New creates the Enviame transport use case
func New(
	client domain.IEnviameClient,
	logger log.ILogger,
) *UseCase {
	return &UseCase{
		client: client,
		log:    logger.WithModule("enviame.usecase"),
	}
}
