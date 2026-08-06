package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const dailyRunHour = 9

type LowBalanceWorker struct {
	uc  ports.IWalletUseCase
	log log.ILogger
}

func NewLowBalanceWorker(uc ports.IWalletUseCase, logger log.ILogger) *LowBalanceWorker {
	return &LowBalanceWorker{uc: uc, log: logger.WithModule("pay.low_balance_worker")}
}

func (w *LowBalanceWorker) Start(ctx context.Context) {
	for {
		next := nextDailyRun(time.Now())
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			w.runCheck(ctx)
		}
	}
}

func nextDailyRun(now time.Time) time.Time {
	loc := now.Location()
	target := time.Date(now.Year(), now.Month(), now.Day(), dailyRunHour, 0, 0, 0, loc)
	if !now.Before(target) {
		target = target.Add(24 * time.Hour)
	}
	return target
}

func (w *LowBalanceWorker) runCheck(ctx context.Context) {
	if err := w.uc.CheckLowBalances(ctx); err != nil {
		w.log.Error(ctx).Err(err).Msg("failed to check low wallet balances")
	}
}
