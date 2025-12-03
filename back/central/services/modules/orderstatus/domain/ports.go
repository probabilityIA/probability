package domain

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// IRepository define la interfaz para el almacenamiento de mapeos de estado
type IRepository interface {
	Create(ctx context.Context, mapping *models.OrderStatusMapping) error
	GetByID(ctx context.Context, id uint) (*models.OrderStatusMapping, error)
	List(ctx context.Context, filters map[string]interface{}) ([]models.OrderStatusMapping, int64, error)
	Update(ctx context.Context, mapping *models.OrderStatusMapping) error
	Delete(ctx context.Context, id uint) error
	ToggleActive(ctx context.Context, id uint) (*models.OrderStatusMapping, error)
	Exists(ctx context.Context, integrationType, originalStatus string) (bool, error)
}
