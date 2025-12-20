package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/fulfillmentstatus/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type FulfillmentStatusHandlers struct {
	uc           app.IUseCase
	logger       log.ILogger
	imageURLBase string
}

func New(uc app.IUseCase, logger log.ILogger, environment env.IConfig) *FulfillmentStatusHandlers {
	imageURLBase := environment.Get("IMAGE_URL_BASE")
	if imageURLBase == "" {
		imageURLBase = "https://storage.googleapis.com"
	}

	return &FulfillmentStatusHandlers{
		uc:           uc,
		logger:       logger,
		imageURLBase: imageURLBase,
	}
}

func (h *FulfillmentStatusHandlers) getImageURLBase() string {
	return h.imageURLBase
}
