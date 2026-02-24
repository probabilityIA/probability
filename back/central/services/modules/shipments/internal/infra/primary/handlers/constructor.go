package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// Handlers contiene todos los handlers del módulo shipments
type Handlers struct {
	uc              *usecases.UseCases
	transportPub    domain.ITransportRequestPublisher // Async: quote, generate, track, cancel
	carrierResolver domain.ICarrierResolver           // Resolves active shipping carrier per business
}

// New crea una nueva instancia de Handlers
func New(uc *usecases.UseCases, transportPub domain.ITransportRequestPublisher, carrierResolver domain.ICarrierResolver) *Handlers {
	return &Handlers{
		uc:              uc,
		transportPub:    transportPub,
		carrierResolver: carrierResolver,
	}
}
