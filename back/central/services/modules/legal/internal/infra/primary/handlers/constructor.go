package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handler struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(useCase app.IUseCase, logger log.ILogger) *Handler {
	return &Handler{uc: useCase, log: logger}
}
