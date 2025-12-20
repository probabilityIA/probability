package domain

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// IRepository define la interfaz para el almacenamiento de estados de fulfillment
type IRepository interface {
	// GetFulfillmentStatusByCode obtiene un estado de fulfillment por su código
	GetFulfillmentStatusByCode(ctx context.Context, code string) (*models.FulfillmentStatus, error)
	// GetFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
	GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error)
	// ListFulfillmentStatuses lista todos los estados de fulfillment
	ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]models.FulfillmentStatus, error)
}
