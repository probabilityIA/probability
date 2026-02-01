package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

// IRepository define la interfaz para el almacenamiento de mapeos de estado
// PURO - Solo tipos de dominio, sin tipos de infraestructura
type IRepository interface {
	Create(ctx context.Context, mapping *entities.OrderStatusMapping) error
	GetByID(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error)
	Update(ctx context.Context, mapping *entities.OrderStatusMapping) error
	Delete(ctx context.Context, id uint) error
	ToggleActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	Exists(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error)
	GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error)
}
