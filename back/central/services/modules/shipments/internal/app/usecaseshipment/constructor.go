package usecaseshipment

import (
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// UseCaseShipment contiene los casos de uso CRUD básicos de envíos
type UseCaseShipment struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCaseShipment
func New(repo domain.IRepository) *UseCaseShipment {
	return &UseCaseShipment{
		repo: repo,
	}
}

