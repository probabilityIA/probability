package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD    *usecaseorder.UseCaseOrder
	orderMapping domain.IOrderMappingUseCase
}

// New crea una nueva instancia de Handlers
func New(orderCRUD *usecaseorder.UseCaseOrder, orderMapping domain.IOrderMappingUseCase) *Handlers {
	return &Handlers{
		orderCRUD:    orderCRUD,
		orderMapping: orderMapping,
	}
}
