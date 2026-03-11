package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	GetConfig(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error)
	UpdateConfig(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error)
}

type UseCase struct {
	repo   ports.IRepository
	logger log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, logger: logger}
}
