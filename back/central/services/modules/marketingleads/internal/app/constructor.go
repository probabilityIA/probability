package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	CreateLead(ctx context.Context, dto dtos.CreateLeadDTO) error
}

type UseCase struct {
	repo   ports.IRepository
	wa     ports.IWhatsAppSender
	logger log.ILogger
}

func New(repo ports.IRepository, wa ports.IWhatsAppSender, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, wa: wa, logger: logger}
}
