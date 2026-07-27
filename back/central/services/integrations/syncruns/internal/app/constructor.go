package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	repo   domain.IRepository
	logger log.ILogger
}

func New(repo domain.IRepository, logger log.ILogger) domain.IUseCase {
	return &useCase{
		repo:   repo,
		logger: logger.WithModule("syncruns"),
	}
}
