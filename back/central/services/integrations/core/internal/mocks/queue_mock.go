package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// QueueMock implementa rabbitmq.IQueue para testing
type QueueMock struct {
	mock.Mock
}

func (m *QueueMock) Publish(ctx context.Context, queueName string, message []byte) error {
	args := m.Called(ctx, queueName, message)
	return args.Error(0)
}

func (m *QueueMock) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error {
	args := m.Called(ctx, exchangeName, routingKey, message)
	return args.Error(0)
}

func (m *QueueMock) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	args := m.Called(ctx, queueName, handler)
	return args.Error(0)
}

func (m *QueueMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *QueueMock) DeclareQueue(queueName string, durable bool) error {
	args := m.Called(queueName, durable)
	return args.Error(0)
}

func (m *QueueMock) DeclareExchange(exchangeName string, exchangeType string, durable bool) error {
	args := m.Called(exchangeName, exchangeType, durable)
	return args.Error(0)
}

func (m *QueueMock) BindQueue(queueName string, exchangeName string, routingKey string) error {
	args := m.Called(queueName, exchangeName, routingKey)
	return args.Error(0)
}

func (m *QueueMock) Ping() error {
	args := m.Called()
	return args.Error(0)
}
