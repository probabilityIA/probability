package app

import (
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	repo     ports.IRepository
	channels ports.IChannelRegistry
	logger   log.ILogger
}

func New(repo ports.IRepository, channels ports.IChannelRegistry, logger log.ILogger) ports.IUseCase {
	return &useCase{
		repo:     repo,
		channels: channels,
		logger:   logger.WithModule("orderscompare"),
	}
}
