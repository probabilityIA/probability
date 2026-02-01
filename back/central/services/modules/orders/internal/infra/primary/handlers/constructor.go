package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
)

// Handlers contiene todos los handlers del módulo orders
type Handlers struct {
	orderCRUD             ports.IOrderUseCase // Interfaz en lugar de tipo concreto
	orderMapping          ports.IOrderMappingUseCase
	requestConfirmationUC usecaseorder.IRequestConfirmationUseCase
}

// New crea una nueva instancia de Handlers
func New(
	orderCRUD ports.IOrderUseCase, // Interfaz en lugar de tipo concreto
	orderMapping ports.IOrderMappingUseCase,
	requestConfirmationUC usecaseorder.IRequestConfirmationUseCase,
) *Handlers {
	return &Handlers{
		orderCRUD:             orderCRUD,
		orderMapping:          orderMapping,
		requestConfirmationUC: requestConfirmationUC,
	}
}
