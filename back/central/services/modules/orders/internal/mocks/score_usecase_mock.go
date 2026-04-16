package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

// ScoreUseCaseMock es un mock del caso de uso de score
type ScoreUseCaseMock struct {
	mock.Mock
}

func (m *ScoreUseCaseMock) CalculateOrderScore(order *entities.ProbabilityOrder) (float64, []string) {
	args := m.Called(order)
	return args.Get(0).(float64), args.Get(1).([]string)
}

func (m *ScoreUseCaseMock) CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}
