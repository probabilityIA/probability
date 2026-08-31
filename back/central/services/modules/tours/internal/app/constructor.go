package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	ListProgress(ctx context.Context, input dtos.ListProgressInput) (*dtos.ListProgressResult, error)
	SaveProgress(ctx context.Context, input dtos.SaveProgressInput) (*entities.TourProgress, error)
	ResetTour(ctx context.Context, input dtos.ResetInput) error
	ResetAll(ctx context.Context, input dtos.ListProgressInput) error
	SkipAll(ctx context.Context, input dtos.SkipAllInput) error
}

type useCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &useCase{repo: repo, log: logger}
}
