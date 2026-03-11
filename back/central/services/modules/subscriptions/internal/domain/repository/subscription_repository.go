package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

type SubscriptionRepository interface {
	// Obtiene la suscripción actual o más reciente de un negocio
	GetLatestByBusinessID(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error)

	// Lista el historial de suscripciones de un negocio
	ListByBusinessID(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error)

	// Crea un nuevo registro de suscripción/pago
	Create(ctx context.Context, subscription *entities.BusinessSubscription) error

	// Actualiza el estado de suscripción y fecha de fin en el modelo Business
	UpdateBusinessSubscriptionStatus(ctx context.Context, businessID uint, status string, endDate *string) error
}
