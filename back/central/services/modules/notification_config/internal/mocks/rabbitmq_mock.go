package mocks

import "context"

// RabbitMQMock - Mock de rabbitmq.IQueue para testing
type RabbitMQMock struct {
	PublishFn           func(ctx context.Context, queueName string, message []byte) error
	PublishToExchangeFn func(ctx context.Context, exchangeName string, routingKey string, message []byte) error
	ConsumeFn           func(ctx context.Context, queueName string, handler func([]byte) error) error
	CloseFn             func() error
	DeclareQueueFn      func(queueName string, durable bool) error
	DeclareExchangeFn   func(exchangeName string, exchangeType string, durable bool) error
	BindQueueFn         func(queueName string, exchangeName string, routingKey string) error
	PingFn              func() error
}

func (m *RabbitMQMock) Publish(ctx context.Context, queueName string, message []byte) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, queueName, message)
	}
	return nil
}

func (m *RabbitMQMock) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error {
	if m.PublishToExchangeFn != nil {
		return m.PublishToExchangeFn(ctx, exchangeName, routingKey, message)
	}
	return nil
}

func (m *RabbitMQMock) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	if m.ConsumeFn != nil {
		return m.ConsumeFn(ctx, queueName, handler)
	}
	return nil
}

func (m *RabbitMQMock) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}

func (m *RabbitMQMock) DeclareQueue(queueName string, durable bool) error {
	if m.DeclareQueueFn != nil {
		return m.DeclareQueueFn(queueName, durable)
	}
	return nil
}

func (m *RabbitMQMock) DeclareExchange(exchangeName string, exchangeType string, durable bool) error {
	if m.DeclareExchangeFn != nil {
		return m.DeclareExchangeFn(exchangeName, exchangeType, durable)
	}
	return nil
}

func (m *RabbitMQMock) BindQueue(queueName string, exchangeName string, routingKey string) error {
	if m.BindQueueFn != nil {
		return m.BindQueueFn(queueName, exchangeName, routingKey)
	}
	return nil
}

func (m *RabbitMQMock) Ping() error {
	if m.PingFn != nil {
		return m.PingFn()
	}
	return nil
}
