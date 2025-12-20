package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type OrderStatusMappingHandlers struct {
	uc     app.IUseCase
	logger log.ILogger
	env    env.IConfig
}

func New(uc app.IUseCase, logger log.ILogger, env env.IConfig) *OrderStatusMappingHandlers {
	return &OrderStatusMappingHandlers{
		uc:     uc,
		logger: logger,
		env:    env,
	}
}

// getImageURLBase obtiene la URL base de S3 para construir URLs completas
func (h *OrderStatusMappingHandlers) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}
