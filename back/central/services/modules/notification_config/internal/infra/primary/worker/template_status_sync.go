package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	templateStatusSyncFirstRun = 3 * time.Minute
	templateStatusSyncInterval = 10 * time.Minute
)

type TemplateStatusSync struct {
	useCase templates.IUseCase
	logger  log.ILogger
}

func NewTemplateStatusSync(useCase templates.IUseCase, logger log.ILogger) *TemplateStatusSync {
	return &TemplateStatusSync{useCase: useCase, logger: logger.WithModule("template_status_sync_worker")}
}

func (w *TemplateStatusSync) Start(ctx context.Context) {
	timer := time.NewTimer(templateStatusSyncFirstRun)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			synced, err := w.useCase.SyncAllPendingStatuses(ctx)
			if err != nil {
				w.logger.Error(ctx).Err(err).Msg("Error consultando las plantillas pendientes")
			} else if synced > 0 {
				w.logger.Info(ctx).Int("negocios", synced).Msg("Estado de plantillas pendientes pedido a Meta")
			}
			timer.Reset(templateStatusSyncInterval)
		}
	}
}
