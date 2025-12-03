package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type OrderStatusMappingHandlers struct {
	uc     app.IUseCase
	logger log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) *OrderStatusMappingHandlers {
	return &OrderStatusMappingHandlers{
		uc:     uc,
		logger: logger,
	}
}
