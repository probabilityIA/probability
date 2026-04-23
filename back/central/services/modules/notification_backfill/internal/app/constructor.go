package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	registry          ports.ISelectorRegistry
	store             ports.IJobStore
	progressPublisher ports.IProgressPublisher
	log               log.ILogger
}

func New(
	registry ports.ISelectorRegistry,
	store ports.IJobStore,
	progressPublisher ports.IProgressPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		registry:          registry,
		store:             store,
		progressPublisher: progressPublisher,
		log:               logger.WithModule("notification_backfill.usecase"),
	}
}
