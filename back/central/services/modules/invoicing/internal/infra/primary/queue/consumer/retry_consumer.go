package consumer

import (
	"context"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// RetryConsumer procesa reintentos de facturas fallidas
type RetryConsumer struct {
	repo ports.IRepository
	useCase     ports.IUseCase
	log         log.ILogger
	ticker      *time.Ticker
}

// NewRetryConsumer crea un nuevo retry consumer
func NewRetryConsumer(
	repo ports.IRepository,
	useCase ports.IUseCase,
	logger log.ILogger,
) *RetryConsumer {
	return &RetryConsumer{
		repo: repo,
		useCase:     useCase,
		log:         logger.WithModule("invoicing.retry_consumer"),
	}
}

// Start inicia el procesamiento de reintentos (cada ~5 minutos con jitter)
func (c *RetryConsumer) Start(ctx context.Context) {
	// Agregar jitter aleatorio (0-60s) para evitar thundering herd en múltiples instancias
	jitter := time.Duration(rand.Intn(60)) * time.Second
	interval := 5*time.Minute + jitter
	c.ticker = time.NewTicker(interval)

	// Primera ejecución inmediata
	c.processRetries(ctx)

	// Luego ejecutar cada 5 minutos
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

// Stop detiene el retry consumer
func (c *RetryConsumer) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.log.Info(context.Background()).Msg("Retry consumer stopped")
	}
}

// processRetries procesa todos los reintentos pendientes
func (c *RetryConsumer) processRetries(ctx context.Context) {
	c.log.Debug(ctx).Msg("Processing pending retries")

	// Obtener logs con reintentos pendientes
	// - status = failed
	// - retry_count < max_retries (3)
	// - next_retry_at <= now
	logs, err := c.repo.GetPendingSyncLogRetries(ctx, 50) // Máximo 50 por batch
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to get pending retries")
		return
	}

	if len(logs) == 0 {
		c.log.Debug(ctx).Msg("No pending retries found")
		return
	}

	c.log.Info(ctx).
		Int("count", len(logs)).
		Msg("Found pending retries")

	// Procesar cada reintento
	successCount := 0
	failCount := 0

	for _, syncLog := range logs {
		if err := c.retryInvoice(ctx, syncLog.InvoiceID); err != nil {
			c.log.Error(ctx).
				Err(err).
				Uint("invoice_id", syncLog.InvoiceID).
				Msg("Failed to retry invoice")
			failCount++
		} else {
			successCount++
		}
	}

	c.log.Info(ctx).
		Int("success", successCount).
		Int("failed", failCount).
		Int("total", len(logs)).
		Msg("Retry batch completed")
}

// retryInvoice reintenta una factura específica
func (c *RetryConsumer) retryInvoice(ctx context.Context, invoiceID uint) error {
	c.log.Debug(ctx).
		Uint("invoice_id", invoiceID).
		Msg("Retrying invoice")

	// Usar el caso de uso de reintento
	err := c.useCase.RetryInvoice(ctx, invoiceID)
	if err != nil {
		// Log pero no fallar - el sync log ya registrará el intento
		c.log.Warn(ctx).
			Err(err).
			Uint("invoice_id", invoiceID).
			Msg("Retry invoice failed - will try again later if within max retries")
		return err
	}

	c.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Msg("Invoice retry successful")

	return nil
}

