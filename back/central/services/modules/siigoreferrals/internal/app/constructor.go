package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/siigoreferrals/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateReferral(ctx context.Context, dto dtos.CreateReferralDTO) error
	ListReferrals(ctx context.Context, page, pageSize int) ([]entities.SiigoReferral, dtos.ListReferralsResult, error)
}

type UseCase struct {
	repo   ports.IRepository
	logger log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, logger: logger}
}
