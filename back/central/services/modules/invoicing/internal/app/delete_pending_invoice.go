package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

const minQueryAttemptsForDelete = 3

// DeletePendingInvoice elimina una factura que está en estado "pending" y tiene
// al menos 3 intentos de consulta (query) sin resolverse.
func (uc *useCase) DeletePendingInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Deleting pending invoice")

	// 1. Verificar que la factura existe
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Invoice not found")
		return errors.ErrInvoiceNotFound
	}

	// 2. Verificar que está en estado pending
	if invoice.Status != constants.InvoiceStatusPending {
		uc.log.Warn(ctx).Str("status", invoice.Status).Msg("Invoice is not in pending status")
		return errors.ErrInvoiceNotPending
	}

	// 3. Verificar que tiene al menos 3 intentos de consulta
	syncLogs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get sync logs")
		return fmt.Errorf("failed to get sync logs: %w", err)
	}

	queryAttempts := 0
	for _, log := range syncLogs {
		if log.OperationType == constants.OperationTypeQuery {
			queryAttempts++
		}
	}

	if queryAttempts < minQueryAttemptsForDelete {
		uc.log.Warn(ctx).
			Int("query_attempts", queryAttempts).
			Int("min_required", minQueryAttemptsForDelete).
			Msg("Insufficient query attempts for deletion")
		return errors.ErrInsufficientQueryAttempts
	}

	// 4. Cancelar sync logs pendientes antes de eliminar
	for _, log := range syncLogs {
		if log.Status == constants.SyncStatusPending || log.Status == constants.SyncStatusFailed {
			log.Status = constants.SyncStatusCancelled
			log.NextRetryAt = nil
			if err := uc.repo.UpdateInvoiceSyncLog(ctx, log); err != nil {
				uc.log.Error(ctx).Err(err).Uint("log_id", log.ID).Msg("Failed to cancel sync log")
			}
		}
	}

	// 5. Eliminar la factura (soft delete)
	if err := uc.repo.DeleteInvoice(ctx, invoiceID); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to delete invoice")
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoiceID).
		Int("query_attempts", queryAttempts).
		Msg("Pending invoice deleted successfully")

	return nil
}
