package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCaseAction define la interfaz para los casos de uso de actions
type IUseCaseAction interface {
	GetActions(ctx context.Context, page, pageSize int, name string) (*domain.ActionListDTO, error)
	GetActionByID(ctx context.Context, id uint) (*domain.ActionDTO, error)
	CreateAction(ctx context.Context, action domain.CreateActionDTO) (*domain.ActionDTO, error)
	UpdateAction(ctx context.Context, id uint, action domain.UpdateActionDTO) (*domain.ActionDTO, error)
	DeleteAction(ctx context.Context, id uint) (string, error)
}

type ActionUseCase struct {
	repository domain.IRepository
	logger     log.ILogger
}

// New crea una nueva instancia del caso de uso de actions
func New(repository domain.IRepository, logger log.ILogger) IUseCaseAction {
	return &ActionUseCase{
		repository: repository,
		logger:     logger,
	}
}
