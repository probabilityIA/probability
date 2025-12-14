package usecaseorderscore

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// UseCaseOrderScore contiene la lógica de negocio para calcular el score de las órdenes
type UseCaseOrderScore struct {
	repo domain.IRepository
}

// IOrderScoreUseCase define la interfaz para el caso de uso de cálculo de score
type IOrderScoreUseCase interface {
	CalculateOrderScore(order *domain.ProbabilityOrder) (float64, []string)
	CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error
}

// New crea una nueva instancia de UseCaseOrderScore
func New(repo domain.IRepository) IOrderScoreUseCase {
	return &UseCaseOrderScore{
		repo: repo,
	}
}
