package consumer

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	retryBatchSize = 50
	retryBaseDelay = 15 * time.Minute
	retryMaxDelay  = 6 * time.Hour
)

type retryOutcome int

const (
	outcomeDispatched retryOutcome = iota
	outcomeRetired
	outcomeDeferred
	outcomeHandled
)

var permanentRetryErrors = []error{
	domainerrors.ErrRetryNotAllowed,
	domainerrors.ErrMaxRetriesExceeded,
	domainerrors.ErrInvoiceNotFound,
	domainerrors.ErrProviderNotConfigured,
	domainerrors.ErrSyncLogNotFound,
}

type RetryConsumer struct {
	repo    ports.IRepository
	useCase ports.IUseCase
	log     log.ILogger
	ticker  *time.Ticker
}

func NewRetryConsumer(
	repo ports.IRepository,
	useCase ports.IUseCase,
	logger log.ILogger,
) *RetryConsumer {
	return &RetryConsumer{
		repo:    repo,
		useCase: useCase,
		log:     logger.WithModule("invoicing.retry_consumer"),
	}
}

func (c *RetryConsumer) Start(ctx context.Context) {
	jitter := time.Duration(rand.Intn(60)) * time.Second
	interval := 5*time.Minute + jitter
	c.ticker = time.NewTicker(interval)

	c.processRetries(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info(ctx).Msg("Retry consumer stopped")
				return
			case <-c.ticker.C:
				c.processRetries(ctx)
			}
		}
	}()
}

func (c *RetryConsumer) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.log.Info(context.Background()).Msg("Retry consumer stopped")
	}
}

func (c *RetryConsumer) processRetries(ctx context.Context) {
	logs, err := c.repo.GetPendingSyncLogRetries(ctx, retryBatchSize)
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al obtener reintentos pendientes")
		return
	}

	if len(logs) == 0 {
		return
	}

	var dispatched, retired, deferred, handled int
	for _, syncLog := range logs {
		switch c.processOne(ctx, syncLog) {
		case outcomeDispatched:
			dispatched++
		case outcomeRetired:
			retired++
		case outcomeDeferred:
			deferred++
		case outcomeHandled:
			handled++
		}
	}

	event := c.log.Info(ctx)
	if retired > 0 || deferred > 0 {
		event = c.log.Warn(ctx)
	}
	event.
		Int("total", len(logs)).
		Int("dispatched", dispatched).
		Int("retired", retired).
		Int("deferred", deferred).
		Int("handled", handled).
		Msg("Retry/check batch completed")
}

func (c *RetryConsumer) processOne(ctx context.Context, syncLog *entities.InvoiceSyncLog) retryOutcome {
	invoice, err := c.repo.GetInvoiceByID(ctx, syncLog.InvoiceID)
	if err != nil || invoice == nil {
		c.retire(ctx, syncLog, "la factura referenciada ya no existe")
		return outcomeRetired
	}

	switch invoice.Status {
	case constants.InvoiceStatusPending:
		err = c.useCase.CheckPendingInvoice(ctx, syncLog.InvoiceID)
	case constants.InvoiceStatusFailed:
		err = c.useCase.RetryInvoice(ctx, syncLog.InvoiceID, false)
	default:
		c.retire(ctx, syncLog, "la factura quedo en "+invoice.Status+" y no necesita reintento")
		return outcomeRetired
	}

	if err == nil {
		return outcomeDispatched
	}

	if !c.stillEligible(ctx, syncLog) {
		return outcomeHandled
	}

	if isPermanentRetryError(err) {
		c.retire(ctx, syncLog, err.Error())
		return outcomeRetired
	}

	c.deferRetry(ctx, syncLog, err)
	return outcomeDeferred
}

func (c *RetryConsumer) stillEligible(ctx context.Context, syncLog *entities.InvoiceSyncLog) bool {
	logs, err := c.repo.GetSyncLogsByInvoiceID(ctx, syncLog.InvoiceID)
	if err != nil {
		return false
	}
	for _, l := range logs {
		if l.ID != syncLog.ID {
			continue
		}
		if l.NextRetryAt == nil || l.RetryCount >= l.MaxRetries {
			return false
		}
		return l.Status == constants.SyncStatusFailed || l.Status == constants.SyncStatusPending
	}
	return false
}

func (c *RetryConsumer) retire(ctx context.Context, syncLog *entities.InvoiceSyncLog, reason string) {
	syncLog.Status = constants.SyncStatusCancelled
	syncLog.NextRetryAt = nil
	syncLog.ErrorMessage = &reason

	if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", syncLog.InvoiceID).
			Uint("sync_log_id", syncLog.ID).
			Msg("No se pudo retirar el sync log: vuelve en la proxima vuelta")
		return
	}

	c.log.Warn(ctx).
		Uint("invoice_id", syncLog.InvoiceID).
		Uint("sync_log_id", syncLog.ID).
		Str("reason", reason).
		Msg("Reintento retirado: repetirlo no puede cambiar el resultado")
}

func (c *RetryConsumer) deferRetry(ctx context.Context, syncLog *entities.InvoiceSyncLog, cause error) {
	syncLog.RetryCount++
	nextRetry := time.Now().Add(retryBackoff(syncLog.RetryCount))
	syncLog.NextRetryAt = &nextRetry
	message := cause.Error()
	syncLog.ErrorMessage = &message

	if err := c.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
		c.log.Error(ctx).
			Err(err).
			Uint("invoice_id", syncLog.InvoiceID).
			Uint("sync_log_id", syncLog.ID).
			Msg("No se pudo aplazar el sync log: vuelve en la proxima vuelta")
		return
	}

	c.log.Warn(ctx).
		Err(cause).
		Uint("invoice_id", syncLog.InvoiceID).
		Uint("sync_log_id", syncLog.ID).
		Int("retry_count", syncLog.RetryCount).
		Int("max_retries", syncLog.MaxRetries).
		Time("next_retry_at", nextRetry).
		Msg("Reintento aplazado por error transitorio")
}

func retryBackoff(attempt int) time.Duration {
	delay := retryBaseDelay
	for i := 1; i < attempt; i++ {
		if delay >= retryMaxDelay {
			break
		}
		delay *= 2
	}
	if delay > retryMaxDelay {
		return retryMaxDelay
	}
	return delay
}

func isPermanentRetryError(err error) bool {
	for _, target := range permanentRetryErrors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
