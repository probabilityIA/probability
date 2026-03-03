package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	emailClient ports.IEmailClient
	resultPub   ports.IResultPublisher
	logger      log.ILogger
}

// New crea un nuevo caso de uso de email
func New(client ports.IEmailClient, resultPub ports.IResultPublisher, logger log.ILogger) ports.IEmailUseCase {
	return &useCase{
		emailClient: client,
		resultPub:   resultPub,
		logger:      logger,
	}
}
