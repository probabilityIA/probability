package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateSprintDTO) (*entities.Sprint, error)
	Get(ctx context.Context, id uint) (*entities.Sprint, error)
	List(ctx context.Context, params dtos.ListSprintsParams) ([]entities.Sprint, int64, error)
	Update(ctx context.Context, dto dtos.UpdateSprintDTO) (*entities.Sprint, error)
	Delete(ctx context.Context, id uint) error
	ChangeStatus(ctx context.Context, dto dtos.ChangeSprintStatusDTO) (*entities.Sprint, error)
}

type UseCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, log: logger}
}
