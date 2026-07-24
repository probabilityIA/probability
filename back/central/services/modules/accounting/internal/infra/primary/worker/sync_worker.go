package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

const syncInterval = 15 * time.Minute

type SyncWorker struct {
	uc  app.IUseCase
	log log.ILogger
}

func NewSyncWorker(uc app.IUseCase, logger log.ILogger) *SyncWorker {
	return &SyncWorker{uc: uc, log: logger.WithModule("accounting.sync")}
}

func (w *SyncWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *SyncWorker) runOnce(ctx context.Context) {
	result, err := w.uc.Sync(ctx)
	if err != nil {
		w.log.Error(ctx).Err(err).Msg("Accounting sync worker: cycle failed")
		return
	}
	if result.Created > 0 {
		w.log.Info(ctx).Int("created", result.Created).Interface("by_source", result.BySource).Msg("Accounting sync worker: cycle complete")
	}
}
