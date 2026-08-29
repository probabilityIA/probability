package retention

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	DefaultInterval  = 24 * time.Hour
	DefaultBatchSize = 1000
	DefaultMaxBatch  = 50
)

type Task struct {
	Name     string
	Days     int
	Interval time.Duration
	Run      func(ctx context.Context, days, batchSize int) (int64, error)
}

func Start(ctx context.Context, logger log.ILogger, task Task) {
	if task.Run == nil || task.Days <= 0 {
		return
	}
	interval := task.Interval
	if interval <= 0 {
		interval = DefaultInterval
	}

	run := func() {
		var total int64
		for i := 0; i < DefaultMaxBatch; i++ {
			affected, err := task.Run(ctx, task.Days, DefaultBatchSize)
			if err != nil {
				logger.Warn(ctx).Err(err).Str("task", task.Name).Msg("retention cleanup failed")
				return
			}
			total += affected
			if affected < int64(DefaultBatchSize) {
				break
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
		if total > 0 {
			logger.Info(ctx).Int64("affected", total).Int("retention_days", task.Days).
				Str("task", task.Name).Msg("retention cleanup done")
		}
	}

	run()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
