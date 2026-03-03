package handlers

import (
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type IHandler interface {
	RegisterRoutes(group interface{})
}

type Handlers struct {
	useCase ports.IUseCase
	log     log.ILogger
}

func New(useCase ports.IUseCase, logger log.ILogger) *Handlers {
	return &Handlers{
		useCase: useCase,
		log:     logger,
	}
}
