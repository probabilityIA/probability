package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD    *usecaseorder.UseCaseOrder
	orderMapping usecaseordermapping.IOrderMappingUseCase
}

// New crea una nueva instancia de Handlers
func New(orderCRUD *usecaseorder.UseCaseOrder, orderMapping usecaseordermapping.IOrderMappingUseCase) *Handlers {
	return &Handlers{
		orderCRUD:    orderCRUD,
		orderMapping: orderMapping,
	}
}
