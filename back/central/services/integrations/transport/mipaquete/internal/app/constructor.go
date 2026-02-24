package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCase implements the transport operations for MiPaquete
type UseCase struct {
	client domain.IMiPaqueteClient
	log    log.ILogger
}

// New creates the MiPaquete transport use case
func New(
	client domain.IMiPaqueteClient,
	logger log.ILogger,
) *UseCase {
	return &UseCase{
		client: client,
		log:    logger.WithModule("mipaquete.usecase"),
	}
}
