package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/chatretention"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	chatRetentionFirstRun = 2 * time.Minute
	chatRetentionInterval = 24 * time.Hour
)

type ChatRetention struct {
	useCase chatretention.IUseCase
	logger  log.ILogger
}

func NewChatRetention(useCase chatretention.IUseCase, logger log.ILogger) *ChatRetention {
	return &ChatRetention{useCase: useCase, logger: logger.WithModule("chat_retention_worker")}
}

func (w *ChatRetention) Start(ctx context.Context) {
	timer := time.NewTimer(chatRetentionFirstRun)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if _, err := w.useCase.Run(ctx); err != nil {
				w.logger.Error(ctx).Err(err).Msg("Error borrando el historial de chats vencido")
			}
			timer.Reset(chatRetentionInterval)
		}
	}
}
