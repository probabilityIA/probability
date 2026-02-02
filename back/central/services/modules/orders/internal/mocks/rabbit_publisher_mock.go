package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

// RabbitPublisherMock es un mock del publicador de eventos a RabbitMQ
type RabbitPublisherMock struct {
	mock.Mock
}

func (m *RabbitPublisherMock) PublishOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RabbitPublisherMock) PublishOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RabbitPublisherMock) PublishOrderCancelled(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RabbitPublisherMock) PublishOrderStatusChanged(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error {
	args := m.Called(ctx, order, previousStatus, currentStatus)
	return args.Error(0)
}

func (m *RabbitPublisherMock) PublishConfirmationRequested(ctx context.Context, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *RabbitPublisherMock) PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
	args := m.Called(ctx, event, order)
	return args.Error(0)
}
