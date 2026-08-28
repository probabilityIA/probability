package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	GetPendingDocuments(ctx context.Context, userID uint) (*dtos.PendingDocumentsResult, error)
	AcceptDocuments(ctx context.Context, input dtos.AcceptDocumentsInput) (*dtos.AcceptDocumentsResult, error)
}

type useCase struct {
	repo ports.IRepository
	log  log.ILogger
}

func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &useCase{repo: repo, log: logger}
}
