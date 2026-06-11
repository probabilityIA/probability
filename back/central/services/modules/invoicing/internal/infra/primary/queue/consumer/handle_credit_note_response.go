package consumer

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (c *ResponseConsumer) handleCreditNoteResponse(
	ctx context.Context,
	invoice *entities.Invoice,
	syncLog *entities.InvoiceSyncLog,
	response *dtos.InvoiceResponseMessage,
) {
	note := c.latestPendingCreditNote(ctx, invoice.ID)

	if response.Status != dtos.ResponseStatusSuccess {
		c.log.Error(ctx).
			Uint("invoice_id", invoice.ID).
			Str("error", response.Error).
			Msg("Credit note failed at provider")
		if note != nil {
			note.Status = constants.CreditNoteStatusFailed
			_ = c.repo.UpdateCreditNote(ctx, note)
		}
		c.failSyncLog(ctx, syncLog, response.Error)
		return
	}

	if note != nil {
		note.CreditNoteNumber = response.InvoiceNumber
		if response.ExternalID != "" {
			extID := response.ExternalID
			note.ExternalID = &extID
		}
		if response.CUFE != "" {
			cufe := response.CUFE
			note.CUFE = &cufe
		}
		now := time.Now()
		note.IssuedAt = &now
		note.Status = constants.CreditNoteStatusIssued
		if err := c.repo.UpdateCreditNote(ctx, note); err != nil {
			c.log.Error(ctx).Err(err).Msg("Failed to update credit note")
		}
	}

	if syncLog != nil {
		completedAt := time.Now()
		duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusSuccess
		syncLog.CompletedAt = &completedAt
		syncLog.Duration = &duration
		_ = c.repo.UpdateInvoiceSyncLog(ctx, syncLog)
	}

	_ = rabbitmq.PublishEvent(ctx, c.queue, rabbitmq.EventEnvelope{
		Type:       "invoice.credit_note_created",
		Category:   "invoice",
		BusinessID: invoice.BusinessID,
		Data: map[string]interface{}{
			"invoice_id":         invoice.ID,
			"order_id":           invoice.OrderID,
			"credit_note_number": response.InvoiceNumber,
		},
	})

	c.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("credit_note_number", response.InvoiceNumber).
		Msg("Credit note processed successfully")
}

func (c *ResponseConsumer) latestPendingCreditNote(ctx context.Context, invoiceID uint) *entities.CreditNote {
	notes, err := c.repo.ListCreditNotes(ctx, map[string]interface{}{
		"invoice_id": invoiceID,
		"status":     constants.CreditNoteStatusPending,
		"limit":      1,
	})
	if err != nil || len(notes) == 0 {
		return nil
	}
	return notes[0]
}

func (c *ResponseConsumer) failSyncLog(ctx context.Context, syncLog *entities.InvoiceSyncLog, errMsg string) {
	if syncLog == nil {
		return
	}
	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusFailed
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	if errMsg != "" {
		syncLog.ErrorMessage = &errMsg
	}
	_ = c.repo.UpdateInvoiceSyncLog(ctx, syncLog)
}
