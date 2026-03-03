package app

import (
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type useCase struct {
	repo          ports.IRepository
	centralClient ports.ICentralClient
	log           log.ILogger
}

func New(repo ports.IRepository, centralClient ports.ICentralClient, logger log.ILogger) ports.IUseCase {
	return &useCase{
		repo:          repo,
		centralClient: centralClient,
		log:           logger,
	}
}
