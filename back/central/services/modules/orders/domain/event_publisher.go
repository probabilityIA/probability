package domain

import (
	"context"
)

// IOrderEventPublisher define la interfaz para publicar eventos de órdenes
type IOrderEventPublisher interface {
	// PublishOrderEvent publica un evento de orden a Redis
	PublishOrderEvent(ctx context.Context, event *OrderEvent) error
}
