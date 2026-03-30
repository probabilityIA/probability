package usecaseshipment

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// UseCaseShipment contiene los casos de uso CRUD básicos de envíos
type UseCaseShipment struct {
	repo domain.IRepository
}

// New crea una nueva instancia de UseCaseShipment
func New(repo domain.IRepository) *UseCaseShipment {
	uc := &UseCaseShipment{
		repo: repo,
	}
	
	// One-time migration: Ensure all businesses are active (paid)
	// Passing Background context since this is a startup task
	uc.repo.EnsureAllBusinessesActive(context.Background())

	return uc
}

