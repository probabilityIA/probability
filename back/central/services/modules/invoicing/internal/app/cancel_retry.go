package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// CancelRetry cancela los reintentos pendientes de una factura
func (uc *useCase) CancelRetry(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Cancelling invoice retry")

	// 1. Verificar que la factura existe
	_, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Invoice not found")
		return errors.ErrInvoiceNotFound
	}

	// 2. Obtener logs de sincronización de esta factura
	syncLogs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get sync logs")
		return fmt.Errorf("failed to get sync logs: %w", err)
	}

	if len(syncLogs) == 0 {
		uc.log.Warn(ctx).Msg("No sync logs found for this invoice")
		return errors.ErrSyncLogNotFound
	}

	// 3. Cancelar todos los logs pendientes de reintento
	cancelledCount := 0
	for _, log := range syncLogs {
		// Solo cancelar logs que están fallidos y tienen reintentos pendientes
		if log.Status == constants.SyncStatusFailed && log.NextRetryAt != nil {
			log.Status = constants.SyncStatusCancelled
			log.NextRetryAt = nil
			log.RetryCount = log.MaxRetries // Marcar como si alcanzó el máximo

			if err := uc.repo.UpdateInvoiceSyncLog(ctx, log); err != nil {
				uc.log.Error(ctx).Err(err).Uint("log_id", log.ID).Msg("Failed to cancel sync log")
				continue // Continuar con los demás
			}
			cancelledCount++
		}
	}

	if cancelledCount == 0 {
		uc.log.Warn(ctx).Msg("No pending retries found to cancel")
		return errors.ErrNoRetriesToCancel
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Int("cancelled_count", cancelledCount).
		Msg("Invoice retries cancelled successfully")

	return nil
}
