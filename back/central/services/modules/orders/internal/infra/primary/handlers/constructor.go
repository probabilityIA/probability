package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD             *usecaseorder.UseCaseOrder
	orderMapping          domain.IOrderMappingUseCase
	requestConfirmationUC usecaseorder.IRequestConfirmationUseCase
}

// New crea una nueva instancia de Handlers
func New(
	orderCRUD *usecaseorder.UseCaseOrder,
	orderMapping domain.IOrderMappingUseCase,
	requestConfirmationUC usecaseorder.IRequestConfirmationUseCase,
) *Handlers {
	return &Handlers{
		orderCRUD:             orderCRUD,
		orderMapping:          orderMapping,
		requestConfirmationUC: requestConfirmationUC,
	}
}
