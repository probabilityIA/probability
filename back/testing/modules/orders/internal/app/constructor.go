package app

import (
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type useCase struct {
	repo               ports.IRepository
	centralClient      ports.ICentralClient
	log                log.ILogger
	webhookSimulators  map[string]ports.IWebhookSimulator // keyed by integration_type code (e.g. "Shopify")
}

func New(repo ports.IRepository, centralClient ports.ICentralClient, logger log.ILogger, simulators map[string]ports.IWebhookSimulator) ports.IUseCase {
	return &useCase{
		repo:              repo,
		centralClient:     centralClient,
		log:               logger,
		webhookSimulators: simulators,
	}
}
