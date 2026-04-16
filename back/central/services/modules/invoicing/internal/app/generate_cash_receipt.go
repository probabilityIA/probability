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

// GenerateCashReceipt genera un recibo de caja para una factura ya emitida.
// Solo aplica cuando la config tiene send_cash_receipt=true y la factura está en issued.
func (uc *useCase) GenerateCashReceipt(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Generating cash receipt for issued invoice")

	// 1. Obtener factura
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Validar que esté en estado issued
	if invoice.Status != constants.InvoiceStatusIssued {
		return fmt.Errorf("invoice must be in 'issued' status to generate cash receipt, current: %s", invoice.Status)
	}

	// 3. Validar que tenga número de factura (necesario para consultar el documento)
	if invoice.InvoiceNumber == "" {
		return fmt.Errorf("invoice has no invoice_number — cannot generate cash receipt")
	}

	// 4. Obtener datos de la orden para el integration_id
	order, err := uc.repo.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// 5. Obtener configuración de facturación (sin filtrar por enabled — la factura ya fue emitida,
	// solo necesitamos las credenciales y config de cash receipt)
	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		return fmt.Errorf("failed to get invoicing config: %w", err)
	}
	if config == nil {
		config, err = uc.repo.GetAnyConfigByBusiness(ctx, order.BusinessID)
		if err != nil || config == nil {
			return errors.ErrProviderNotConfigured
		}
	}

	// 6. Validar que la config tenga send_cash_receipt habilitado
	sendCashReceipt := false
	if config.InvoiceConfig != nil {
		if v, ok := config.InvoiceConfig["send_cash_receipt"].(bool); ok {
			sendCashReceipt = v
		}
	}
	if !sendCashReceipt {
		return fmt.Errorf("cash receipt is not enabled in invoicing config")
	}

	// 7. Determinar integración de facturación
	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		return errors.ErrProviderNotConfigured
	}

	// 8. Determinar proveedor
	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error resolving provider for cash receipt, defaulting to softpymes")
		provider = dtos.ProviderSoftpymes
	}

	// 9. Crear sync log
	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoiceID,
		OperationType: constants.OperationTypeCashReceipt,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    0,
		RetryCount:    0,
		TriggeredBy:   constants.TriggerManual,
	}
	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create cash receipt sync log")
	}

	// 10. Construir config con invoice_number para que el consumer lo use
	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			invoiceConfigData[k] = v
		}
	}
	invoiceConfigData["invoice_number"] = invoice.InvoiceNumber

	// 11. Construir mensaje
	correlationID := uuid.New().String()
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: invoice.ID,
		Provider:  provider,
		Operation: dtos.OperationCashReceipt,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			OrderID:       invoice.OrderID,
			Config:        invoiceConfigData,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	// 12. Publicar a RabbitMQ
	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish cash receipt request: " + err.Error()
		syncLog.ErrorMessage = &errorMsg
		_ = uc.repo.UpdateInvoiceSyncLog(ctx, syncLog)
		return fmt.Errorf("failed to publish cash receipt request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Str("invoice_number", invoice.InvoiceNumber).
		Msg("Cash receipt request published")

	return nil
}
