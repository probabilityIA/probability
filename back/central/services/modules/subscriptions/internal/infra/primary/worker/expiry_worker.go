package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

const checkInterval = 24 * time.Hour

type ExpiryWorker struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) *ExpiryWorker {
	return &ExpiryWorker{uc: uc, log: logger}
}

func (w *ExpiryWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	w.runCheck(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runCheck(ctx)
		}
	}
}

func (w *ExpiryWorker) runCheck(ctx context.Context) {
	if err := w.uc.CheckExpiringSubscriptions(ctx); err != nil {
		w.log.Error(ctx).Err(err).Msg("failed to check expiring subscriptions")
	}
}
