package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// EnableRetry re-habilita reintentos automáticos para una factura fallida
func (uc *useCase) EnableRetry(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Enabling invoice retry")

	// 1. Verificar que la factura existe y está en estado failed
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	if invoice.Status != constants.InvoiceStatusFailed {
		return errors.ErrRetryNotAllowed
	}

	// 2. Obtener logs de sincronización
	logs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
	if err != nil || len(logs) == 0 {
		return fmt.Errorf("no sync logs found for invoice")
	}

	// 3. Buscar el último log cancelado y reactivarlo
	// NOTA: No reseteamos retry_count a 0, sino que incrementamos MaxRetries
	// para dar más intentos desde donde quedó (evita reintentos infinitos).
	reenabledCount := 0
	for _, log := range logs {
		if log.Status == constants.SyncStatusCancelled {
			log.Status = constants.SyncStatusFailed
			log.MaxRetries = log.RetryCount + constants.MaxRetries // Dar 3 intentos más desde donde quedó
			nextRetry := time.Now().Add(time.Duration(constants.DefaultRetryIntervalMin) * time.Minute)
			log.NextRetryAt = &nextRetry

			if err := uc.repo.UpdateInvoiceSyncLog(ctx, log); err != nil {
				uc.log.Error(ctx).Err(err).Uint("log_id", log.ID).Msg("Failed to re-enable sync log")
				continue
			}
			reenabledCount++
			break // Solo reactivar el último
		}
	}

	if reenabledCount == 0 {
		return fmt.Errorf("no cancelled retries found to re-enable")
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Int("reenabled_count", reenabledCount).
		Msg("Invoice retries re-enabled successfully")

	return nil
}
