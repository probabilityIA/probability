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

// CancelInvoice anula una factura emitida en el proveedor (Softpymes)
func (uc *useCase) CancelInvoice(ctx context.Context, dto *dtos.CancelInvoiceDTO) error {
	uc.log.Info(ctx).Uint("invoice_id", dto.InvoiceID).Msg("Cancelling invoice")

	// 1. Obtener factura
	invoice, err := uc.repo.GetInvoiceByID(ctx, dto.InvoiceID)
	if err != nil || invoice == nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Solo se pueden anular facturas emitidas
	if invoice.Status == constants.InvoiceStatusCancelled {
		return errors.ErrInvoiceAlreadyCancelled
	}
	if invoice.Status != constants.InvoiceStatusIssued {
		return errors.ErrInvoiceCannotBeCancelled
	}

	// 3. Necesitamos el ExternalID para llamar a Softpymes
	if invoice.ExternalID == nil || *invoice.ExternalID == "" {
		return errors.ErrInvoiceCannotBeCancelled
	}

	// 4. Obtener integración a partir de la factura
	var integrationID uint
	if invoice.InvoicingIntegrationID != nil {
		integrationID = *invoice.InvoicingIntegrationID
	} else if invoice.InvoicingProviderID != nil {
		integrationID = *invoice.InvoicingProviderID
	} else {
		return errors.ErrProviderNotConfigured
	}

	// 5. Determinar proveedor
	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Error al resolver proveedor para cancelación")
		provider = dtos.ProviderSoftpymes
	}

	// 6. Crear sync log de cancelación
	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeCancel,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    0, // La cancelación no tiene reintentos automáticos
		RetryCount:    0,
		TriggeredBy:   constants.TriggerManual,
	}

	if dto.CancelledByUserID > 0 {
		syncLog.UserID = &dto.CancelledByUserID
	}

	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create cancel sync log")
		// Continuamos aunque falle el log
	}

	// 7. Construir config para el consumer de cancelación
	cancelConfig := map[string]interface{}{
		"external_id":   *invoice.ExternalID,
		"cancel_reason": dto.Reason,
	}

	// Obtener config de facturación para base_url / base_url_test / referer
	config, configErr := uc.repo.GetConfigByIntegration(ctx, integrationID)
	if configErr == nil && config != nil {
		cancelConfig["is_testing"] = config.IsTesting
		cancelConfig["base_url"] = config.BaseURL
		cancelConfig["base_url_test"] = config.BaseURLTest
		if config.InvoiceConfig != nil {
			for k, v := range config.InvoiceConfig {
				if _, exists := cancelConfig[k]; !exists {
					cancelConfig[k] = v
				}
			}
		}
	}

	// 8. Generar correlation ID
	correlationID := uuid.New().String()

	// 9. Publicar request de cancelación a RabbitMQ
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: invoice.ID,
		Provider:  provider,
		Operation: dtos.OperationCancel,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			OrderID:       invoice.OrderID,
			Config:        cancelConfig,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to publish cancel request")

		// Marcar sync log como failed
		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errMsg := fmt.Sprintf("Failed to publish cancel request: %s", err.Error())
		syncLog.ErrorMessage = &errMsg
		_ = uc.repo.UpdateInvoiceSyncLog(ctx, syncLog)

		return fmt.Errorf("failed to publish cancel request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("external_id", *invoice.ExternalID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("Cancel request published - waiting for provider response")

	return nil
}
