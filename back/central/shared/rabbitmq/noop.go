package rabbitmq

import (
	"context"
	"errors"
)

// ErrQueueUnavailable se retorna cuando el broker no esta disponible y el
// backend corre en modo degradado.
var ErrQueueUnavailable = errors.New("rabbitmq no disponible")

// noopQueue implementa IQueue sin broker detras. Se usa cuando la conexion a
// RabbitMQ falla al arrancar: permite que los modulos se inicialicen y expongan
// sus rutas HTTP en vez de morir con un nil pointer dereference.
//
// Publicar retorna ErrQueueUnavailable para que el llamador se entere; consumir
// y declarar no hacen nada, porque no hay broker donde registrar nada.
type noopQueue struct{}

// NewNoop construye una IQueue inerte para el modo degradado.
func NewNoop() IQueue {
	return &noopQueue{}
}

func (q *noopQueue) Publish(ctx context.Context, queueName string, message []byte) error {
	return ErrQueueUnavailable
}

func (q *noopQueue) PublishToExchange(ctx context.Context, exchangeName string, routingKey string, message []byte) error {
	return ErrQueueUnavailable
}

func (q *noopQueue) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	return nil
}

func (q *noopQueue) ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	return nil
}

func (q *noopQueue) Close() error {
	return nil
}

func (q *noopQueue) DeclareQueue(queueName string, durable bool) error {
	return nil
}

func (q *noopQueue) DeclareExchange(exchangeName string, exchangeType string, durable bool) error {
	return nil
}

func (q *noopQueue) BindQueue(queueName string, exchangeName string, routingKey string) error {
	return nil
}

func (q *noopQueue) Ping() error {
	return ErrQueueUnavailable
}
