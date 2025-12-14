package usecaseorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// UseCaseOrder contiene los casos de uso CRUD básicos de órdenes
type UseCaseOrder struct {
	repo           domain.IRepository
	eventPublisher domain.IOrderEventPublisher
}

// New crea una nueva instancia de UseCaseOrder
func New(repo domain.IRepository, eventPublisher domain.IOrderEventPublisher) *UseCaseOrder {
	return &UseCaseOrder{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}
