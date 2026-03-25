package app

import (
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseScore struct {
	repo      ports.IRepository
	publisher ports.IScoreEventPublisher
	log       log.ILogger
}

func New(repo ports.IRepository, publisher ports.IScoreEventPublisher, logger log.ILogger) ports.IScoreUseCase {
	return &UseCaseScore{
		repo:      repo,
		publisher: publisher,
		log:       logger,
	}
}
