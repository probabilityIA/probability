package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type PaymentStatusHandlers struct {
	uc           usecases.IUseCase
	logger       log.ILogger
	imageURLBase string
}

func New(uc usecases.IUseCase, logger log.ILogger, environment env.IConfig) *PaymentStatusHandlers {
	imageURLBase := environment.Get("IMAGE_URL_BASE")
	if imageURLBase == "" {
		imageURLBase = "https://storage.googleapis.com"
	}

	return &PaymentStatusHandlers{
		uc:           uc,
		logger:       logger.WithModule("handlers"),
		imageURLBase: imageURLBase,
	}
}

func (h *PaymentStatusHandlers) getImageURLBase() string {
	return h.imageURLBase
}
