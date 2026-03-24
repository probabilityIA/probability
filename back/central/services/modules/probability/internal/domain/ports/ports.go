package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
)

// IRepository define los metodos de acceso a datos para el modulo de probability
type IRepository interface {
	GetOrderForScoring(ctx context.Context, orderID string) (*entities.ScoreOrder, error)
	CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error)
	UpdateOrderScore(ctx context.Context, orderID string, score float64, factors []byte, breakdown []byte) error
	GetCustomerOrderHistory(ctx context.Context, customerID uint, excludeOrderID string) (*entities.CustomerHistory, error)
	GetCustomerDeliveryHistory(ctx context.Context, customerID uint) (*entities.DeliveryHistory, error)
	GetOrderItemCount(ctx context.Context, orderID string) (int, error)
	GetPaymentMethodCategory(ctx context.Context, paymentMethodID uint) (string, error)
}

// IScoreUseCase define los casos de uso del modulo de probability
type IScoreUseCase interface {
	CalculateOrderScore(order *entities.ScoreOrder) (float64, []string, *entities.ScoreBreakdown)
	CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error
}

// IScoreEventPublisher publica eventos de score calculado
type IScoreEventPublisher interface {
	PublishScoreCalculated(ctx context.Context, orderID string, orderNumber string, businessID uint, integrationID uint) error
}
