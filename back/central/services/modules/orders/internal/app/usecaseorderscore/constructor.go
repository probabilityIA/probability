package usecaseorderscore

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
)

// UseCaseOrderScore contiene la lógica de negocio para calcular el score de las órdenes
type UseCaseOrderScore struct {
	repo ports.IRepository
}

// New crea una nueva instancia de UseCaseOrderScore
func New(repo ports.IRepository) ports.IOrderScoreUseCase {
	return &UseCaseOrderScore{
		repo: repo,
	}
}
