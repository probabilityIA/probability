package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/scheduled"
	"github.com/secamc93/probability/back/central/shared/log"
)

const tickInterval = 5 * time.Minute

type Scheduler struct {
	useCase scheduled.IUseCase
	logger  log.ILogger
}

func NewScheduler(useCase scheduled.IUseCase, logger log.ILogger) *Scheduler {
	return &Scheduler{
		useCase: useCase,
		logger:  logger.WithModule("scheduled_notifications_worker"),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	s.logger.Info(ctx).Dur("interval", tickInterval).
		Msg("Worker de notificaciones programadas iniciado")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info(ctx).Msg("Worker de notificaciones programadas detenido")
			return
		case <-ticker.C:
			if err := s.useCase.RunDueRules(ctx); err != nil {
				s.logger.Error(ctx).Err(err).
					Msg("Error en la corrida de reglas programadas")
			}
		}
	}
}
