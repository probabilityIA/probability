package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func (uc *useCase) CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error) {
	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Msg("Creating credit note")

	invoice, err := uc.repo.GetInvoiceByID(ctx, dto.InvoiceID)
	if err != nil || invoice == nil {
		return nil, errors.ErrInvoiceNotFound
	}

	if invoice.Status != constants.InvoiceStatusIssued {
		return nil, fmt.Errorf("solo se pueden crear notas de credito sobre facturas emitidas (estado actual: %s)", invoice.Status)
	}

	if invoice.ExternalID == nil || *invoice.ExternalID == "" {
		return nil, fmt.Errorf("la factura no tiene identificador del proveedor para emitir la nota de credito")
	}

	var integrationID uint
	if invoice.InvoicingIntegrationID != nil {
		integrationID = *invoice.InvoicingIntegrationID
	} else if invoice.InvoicingProviderID != nil {
		integrationID = *invoice.InvoicingProviderID
	} else {
		return nil, errors.ErrProviderNotConfigured
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		provider = dtos.ProviderSoftpymes
	}

	noteType := dto.NoteType
	if noteType == "" {
		noteType = constants.CreditNoteTypeFullRefund
	}
	amount := dto.Amount
	if amount <= 0 {
		amount = invoice.TotalAmount
	}

	creditNote := &entities.CreditNote{
		InvoiceID:      invoice.ID,
		BusinessID:     invoice.BusinessID,
		InternalNumber: fmt.Sprintf("NC-%s", time.Now().Format("20060102-150405")),
		NoteType:       noteType,
		Amount:         amount,
		Currency:       invoice.Currency,
		Reason:         dto.Reason,
		Description:    dto.Description,
		Status:         constants.CreditNoteStatusPending,
	}
	if err := uc.repo.CreateCreditNote(ctx, creditNote); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to persist credit note")
		return nil, fmt.Errorf("failed to create credit note: %w", err)
	}

	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeCreditNote,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    0,
		RetryCount:    0,
		TriggeredBy:   constants.TriggerManual,
	}
	if dto.CreatedByUserID > 0 {
		syncLog.UserID = &dto.CreatedByUserID
	}
	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create credit note sync log")
	}

	creditNoteConfig := map[string]interface{}{
		"external_id":    *invoice.ExternalID,
		"invoice_number": invoice.InvoiceNumber,
		"credit_note_id": creditNote.ID,
		"amount":         amount,
		"reason":         dto.Reason,
		"note_type":      noteType,
		"customer_dni":   invoice.CustomerDNI,
	}

	config, configErr := uc.repo.GetConfigByIntegration(ctx, integrationID)
	if configErr == nil && config != nil {
		creditNoteConfig["is_testing"] = config.IsTesting
		creditNoteConfig["base_url"] = config.BaseURL
		creditNoteConfig["base_url_test"] = config.BaseURLTest
		if config.InvoiceConfig != nil {
			for k, v := range config.InvoiceConfig {
				if _, exists := creditNoteConfig[k]; !exists {
					creditNoteConfig[k] = v
				}
			}
		}
	}

	correlationID := uuid.New().String()
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: invoice.ID,
		Provider:  provider,
		Operation: dtos.OperationCreditNote,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			OrderID:       invoice.OrderID,
			Config:        creditNoteConfig,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to publish credit note request")

		creditNote.Status = constants.CreditNoteStatusFailed
		_ = uc.repo.UpdateCreditNote(ctx, creditNote)

		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errMsg := fmt.Sprintf("Failed to publish credit note request: %s", err.Error())
		syncLog.ErrorMessage = &errMsg
		_ = uc.repo.UpdateInvoiceSyncLog(ctx, syncLog)

		return nil, fmt.Errorf("failed to publish credit note request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Uint("credit_note_id", creditNote.ID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("Credit note request published - waiting for provider response")

	return creditNote, nil
}

func (uc *useCase) GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error) {
	return uc.repo.GetCreditNoteByID(ctx, id)
}

func (uc *useCase) ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error) {
	return uc.repo.ListCreditNotes(ctx, filters)
}
