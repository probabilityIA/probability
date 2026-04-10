package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// DashboardHandlers maneja las peticiones HTTP del módulo dashboard
type DashboardHandlers struct {
	uc     domain.IUseCase
	logger log.ILogger
}

// New crea una nueva instancia de los handlers
func New(uc domain.IUseCase, logger log.ILogger) *DashboardHandlers {
	return &DashboardHandlers{
		uc:     uc,
		logger: logger,
	}
}
