package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type IFlowRepository interface {
	Create(ctx context.Context, flow *entities.Flow) error
	Update(ctx context.Context, flow *entities.Flow) error
	GetByID(ctx context.Context, id, businessID uint) (*entities.Flow, error)
	List(ctx context.Context, businessID uint) ([]entities.Flow, error)
	Delete(ctx context.Context, id, businessID uint) error
	ExistsByName(ctx context.Context, businessID uint, name string, excludeID uint) (bool, error)
}
