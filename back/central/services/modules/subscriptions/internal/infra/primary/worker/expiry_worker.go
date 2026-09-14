package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

const checkTimeZone = "America/Bogota"

var checkHours = []int{8, 20}

type ExpiryWorker struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) *ExpiryWorker {
	return &ExpiryWorker{uc: uc, log: logger}
}

func (w *ExpiryWorker) Start(ctx context.Context) {
	location, err := time.LoadLocation(checkTimeZone)
	if err != nil {
		location = time.UTC
	}

	for {
		wait := time.Until(nextRun(time.Now().In(location), location))
		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			w.runCheck(ctx)
		}
	}
}

func nextRun(now time.Time, location *time.Location) time.Time {
	day := now
	for i := 0; i < 2; i++ {
		for _, hour := range checkHours {
			candidate := time.Date(day.Year(), day.Month(), day.Day(), hour, 0, 0, 0, location)
			if candidate.After(now) {
				return candidate
			}
		}
		day = day.AddDate(0, 0, 1)
	}
	return time.Date(day.Year(), day.Month(), day.Day(), checkHours[0], 0, 0, 0, location)
}

func (w *ExpiryWorker) runCheck(ctx context.Context) {
	if err := w.uc.CheckExpiringSubscriptions(ctx); err != nil {
		w.log.Error(ctx).Err(err).Msg("failed to check expiring subscriptions")
	}
}
