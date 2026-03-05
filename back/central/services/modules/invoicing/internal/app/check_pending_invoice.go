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

// CheckPendingInvoice busca en el proveedor si una factura pendiente de validación DIAN
// ya fue procesada. NO re-envía POST — solo consulta documentos existentes.
func (uc *useCase) CheckPendingInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Checking pending invoice status in provider")

	// 1. Obtener factura
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Solo facturas en pending
	if invoice.Status != constants.InvoiceStatusPending {
		return errors.ErrRetryNotAllowed
	}

	// 3. Obtener sync logs para verificar conteo
	logs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
	if err != nil || len(logs) == 0 {
		return fmt.Errorf("no sync logs found for invoice")
	}

	lastLog := logs[0]

	if lastLog.RetryCount >= lastLog.MaxRetries {
		return errors.ErrMaxRetriesExceeded
	}

	// 4. Cancelar checks pendientes anteriores
	for _, l := range logs {
		if l.Status == constants.SyncStatusPending && l.NextRetryAt != nil {
			l.Status = constants.SyncStatusCancelled
			l.NextRetryAt = nil
			if err := uc.repo.UpdateInvoiceSyncLog(ctx, l); err != nil {
				uc.log.Warn(ctx).Err(err).Uint("sync_log_id", l.ID).Msg("Failed to cancel sync log before check")
			}
		}
	}

	// 5. Crear sync log para este check
	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoiceID,
		OperationType: constants.OperationTypeCreate,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    constants.MaxRetries,
		RetryCount:    lastLog.RetryCount + 1,
		TriggeredBy:   constants.TriggerAuto,
	}

	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create check sync log")
	}

	// 6. Obtener config para autenticación con el proveedor
	order, err := uc.repo.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get order for check")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("failed to get order: %w", err))
		return fmt.Errorf("failed to get order: %w", err)
	}

	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured)
		return errors.ErrProviderNotConfigured
	}
	if config == nil {
		config, err = uc.repo.GetEnabledConfigByBusiness(ctx, order.BusinessID)
		if err != nil || config == nil {
			uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured)
			return errors.ErrProviderNotConfigured
		}
	}

	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured)
		return errors.ErrProviderNotConfigured
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		provider = dtos.ProviderSoftpymes
	}

	// 7. Construir config mínimo (solo necesita auth + order_id para buscar)
	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		invoiceConfigData = config.InvoiceConfig
	}
	invoiceConfigData["is_testing"] = config.IsTesting
	invoiceConfigData["base_url"] = config.BaseURL
	invoiceConfigData["base_url_test"] = config.BaseURLTest

	invoiceData := dtos.InvoiceData{
		IntegrationID: integrationID,
		OrderID:       invoice.OrderID,
		OrderNumber:   order.OrderNumber,
		Currency:      invoice.Currency,
		Config:        invoiceConfigData,
	}

	correlationID := uuid.New().String()

	// 8. Publicar check_status (NO retry) — el consumer solo buscará, no creará
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID:     invoice.ID,
		Provider:      provider,
		Operation:     dtos.OperationCheckStatus,
		InvoiceData:   invoiceData,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to publish check_status request")

		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish check_status: " + err.Error()
		syncLog.ErrorMessage = &errorMsg
		_ = uc.repo.UpdateInvoiceSyncLog(ctx, syncLog)

		return fmt.Errorf("failed to publish check_status request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Int("retry_count", syncLog.RetryCount).
		Msg("Check status request published — searching for existing document")

	return nil
}
