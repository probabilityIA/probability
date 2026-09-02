package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
)

type IRepository interface {
	Create(ctx context.Context, sprint *entities.Sprint) (*entities.Sprint, error)
	GetByID(ctx context.Context, id uint) (*entities.Sprint, error)
	List(ctx context.Context, params dtos.ListSprintsParams) ([]entities.Sprint, int64, error)
	Update(ctx context.Context, id uint, updates map[string]any) (*entities.Sprint, error)
	Delete(ctx context.Context, id uint) error
}
