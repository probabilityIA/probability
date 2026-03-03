package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD             ports.IOrderUseCase
	createUC              ports.IOrderCreateUseCase
	requestConfirmationUC ports.IRequestConfirmationUseCase
	logger                log.ILogger
}

// New crea una nueva instancia de Handlers
func New(
	orderCRUD ports.IOrderUseCase,
	createUC ports.IOrderCreateUseCase,
	requestConfirmationUC ports.IRequestConfirmationUseCase,
	logger log.ILogger,
) *Handlers {
	return &Handlers{
		orderCRUD:             orderCRUD,
		createUC:              createUC,
		requestConfirmationUC: requestConfirmationUC,
		logger:                logger,
	}
}
