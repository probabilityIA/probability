package consumer

import (
	"context"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RetryConsumer procesa reintentos de pagos fallidos
type RetryConsumer struct {
	repo    ports.IRepository
	useCase ports.IUseCase
	log     log.ILogger
	ticker  *time.Ticker
}

// NewRetryConsumer crea un nuevo retry consumer
func NewRetryConsumer(
	repo ports.IRepository,
	useCase ports.IUseCase,
	logger log.ILogger,
) *RetryConsumer {
	return &RetryConsumer{
		repo:    repo,
		useCase: useCase,
		log:     logger.WithModule("pay.retry_consumer"),
	}
}

// Start inicia el procesamiento de reintentos (cada ~5 minutos con jitter)
func (c *RetryConsumer) Start(ctx context.Context) {
	jitter := time.Duration(rand.Intn(60)) * time.Second
	interval := 5*time.Minute + jitter
	c.ticker = time.NewTicker(interval)

	// Primera ejecución inmediata
	c.processRetries(ctx)

	for {
		select {
		case <-ctx.Done():
			c.log.Info(ctx).Msg("Pay retry consumer stopped")
			return
		case <-c.ticker.C:
			c.processRetries(ctx)
		}
	}
}

func (c *RetryConsumer) processRetries(ctx context.Context) {
	logs, err := c.repo.GetPendingSyncLogRetries(ctx, 50)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Error getting pending payment retries")
		return
	}

	if len(logs) == 0 {
		return
	}

	c.log.Info(ctx).Int("count", len(logs)).Msg("Found pending payment retries")

	successCount, failCount := 0, 0
	for _, syncLog := range logs {
		if err := c.useCase.RetryPayment(ctx, syncLog.PaymentTransactionID); err != nil {
			c.log.Error(ctx).Err(err).Uint("transaction_id", syncLog.PaymentTransactionID).Msg("Failed to retry payment")
			failCount++
		} else {
			successCount++
		}
	}

	c.log.Info(ctx).Int("success", successCount).Int("failed", failCount).Int("total", len(logs)).Msg("Pay retry batch completed")
}
