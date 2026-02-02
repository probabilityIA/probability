package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

// EventPublisherMock es un mock del publicador de eventos a Redis
type EventPublisherMock struct {
	mock.Mock
}

func (m *EventPublisherMock) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, event, order)
	return args.Error(0)
}
