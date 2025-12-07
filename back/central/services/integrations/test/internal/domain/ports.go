package domain

import (
	"context"
)

// ───────────────────────────────────────────
//
//	ORDER PUBLISHER INTERFACE
//
// ───────────────────────────────────────────

// IOrderPublisher define la interfaz para publicar órdenes canónicas a RabbitMQ
type IOrderPublisher interface {
	// PublishCanonicalOrder publica una orden canónica a la cola de RabbitMQ
	PublishCanonicalOrder(ctx context.Context, order *CanonicalOrderDTO) error
}

// ───────────────────────────────────────────
//
//	ORDER GENERATOR INTERFACE
//
// ───────────────────────────────────────────

// IOrderGenerator define la interfaz para generar órdenes aleatorias
type IOrderGenerator interface {
	// GenerateRandomOrder genera una orden canónica aleatoria
	GenerateRandomOrder(req *GenerateOrderRequest) *CanonicalOrderDTO
}
