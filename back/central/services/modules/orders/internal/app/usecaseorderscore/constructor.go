package usecaseorderscore

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
)

// UseCaseOrderScore contiene la lógica de negocio para calcular el score de las órdenes
type UseCaseOrderScore struct {
	repo ports.IRepository
}

// IOrderScoreUseCase define la interfaz para el caso de uso de cálculo de score
type IOrderScoreUseCase interface {
	CalculateOrderScore(order *entities.ProbabilityOrder) (float64, []string)
	CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error
}

// New crea una nueva instancia de UseCaseOrderScore
func New(repo ports.IRepository) IOrderScoreUseCase {
	return &UseCaseOrderScore{
		repo: repo,
	}
}
