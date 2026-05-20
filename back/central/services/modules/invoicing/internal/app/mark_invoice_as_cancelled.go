package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func (uc *useCase) MarkInvoiceAsCancelled(ctx context.Context, dto *dtos.MarkInvoiceAsCancelledDTO) error {
	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Str("reason", dto.Reason).Msg("Marking invoice as cancelled")

	invoice, err := uc.repo.GetInvoiceByID(ctx, dto.InvoiceID)
	if err != nil || invoice == nil {
		return errors.ErrInvoiceNotFound
	}

	if invoice.Status == constants.InvoiceStatusCancelled {
		return errors.ErrInvoiceAlreadyCancelled
	}

	now := time.Now()
	invoice.Status = constants.InvoiceStatusCancelled
	invoice.CancelledAt = &now

	if err := uc.repo.UpdateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", dto.InvoiceID).Msg("Failed to update invoice status to cancelled")
		return fmt.Errorf("failed to mark invoice as cancelled: %w", err)
	}

	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeCancel,
		Status:        constants.SyncStatusSuccess,
		StartedAt:     now,
		CompletedAt:   &now,
		TriggeredBy:   constants.TriggerManual,
		UserID:        &dto.CancelledByUserID,
	}

	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to create sync log for mark as cancelled")
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("reason", dto.Reason).
		Msg("Invoice marked as cancelled successfully")

	return nil
}
