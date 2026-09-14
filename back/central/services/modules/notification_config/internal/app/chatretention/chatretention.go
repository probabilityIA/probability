package chatretention

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	Run(ctx context.Context) (entities.ChatPurgeResult, error)
}

type useCase struct {
	repo   ports.IChatRetentionRepository
	logger log.ILogger
	now    func() time.Time
}

func New(repo ports.IChatRetentionRepository, logger log.ILogger) IUseCase {
	return &useCase{repo: repo, logger: logger, now: time.Now}
}

func (uc *useCase) Run(ctx context.Context) (entities.ChatPurgeResult, error) {
	cutoff := uc.now().AddDate(0, 0, -entities.ChatRetentionDays)

	result, err := uc.repo.PurgeChatHistory(ctx, cutoff)
	if err != nil {
		return result, err
	}

	uc.logger.Info(ctx).
		Time("cutoff", cutoff).
		Int64("messages", result.Messages).
		Int64("conversations", result.Conversations).
		Int64("reads", result.Reads).
		Msg("Retencion de chats: historial de mas de un ano eliminado")

	return result, nil
}
