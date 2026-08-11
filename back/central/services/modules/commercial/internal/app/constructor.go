package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	ListProspects(ctx context.Context, filters dtos.ListProspectsFilters) ([]entities.Prospect, dtos.ListResult, error)
	GetStats(ctx context.Context) (*entities.ProspectStats, error)
	UpdateSeen(ctx context.Context, id uint, seen bool) error
}

type UseCase struct {
	repo   ports.IRepository
	logger log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, logger: logger}
}
