package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCaseResource define la interfaz para los casos de uso de recursos
type IUseCaseResource interface {
	GetResources(ctx context.Context, filters domain.ResourceFilters) (*domain.ResourceListDTO, error)
	GetResourceByID(ctx context.Context, id uint) (*domain.ResourceDTO, error)
	CreateResource(ctx context.Context, resource domain.CreateResourceDTO) (*domain.ResourceDTO, error)
	UpdateResource(ctx context.Context, id uint, resource domain.UpdateResourceDTO) (*domain.ResourceDTO, error)
	DeleteResource(ctx context.Context, id uint) (string, error)
}

type ResourceUseCase struct {
	repository domain.IRepository
	logger     log.ILogger
}

// New crea una nueva instancia del caso de uso de recursos
func New(repository domain.IRepository, logger log.ILogger) IUseCaseResource {
	return &ResourceUseCase{
		repository: repository,
		logger:     logger,
	}
}
