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

// RetryInvoice reintenta la creación de una factura fallida (in-place, sin eliminar)
func (uc *useCase) RetryInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Retrying invoice creation")

	// 1. Obtener factura existente (NO se elimina)
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Validar que esté en estado failed (solo se puede reintentar desde failed)
	if invoice.Status != constants.InvoiceStatusFailed {
		return errors.ErrRetryNotAllowed
	}

	// 2.5. Lock optimista: marcar invoice como pending ANTES de llamar al proveedor.
	invoice.Status = constants.InvoiceStatusPending
	if err := uc.repo.UpdateInvoice(ctx, invoice); err != nil {
		return fmt.Errorf("failed to lock invoice for retry: %w", err)
	}

	// 3. Obtener logs de sincronización para verificar reintentos
	logs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
	if err != nil || len(logs) == 0 {
		return fmt.Errorf("no sync logs found for invoice")
	}

	lastLog := logs[0] // Ordenados por created_at DESC

	// 4. Validar que no se haya excedido el máximo de reintentos
	if lastLog.RetryCount >= lastLog.MaxRetries {
		return errors.ErrMaxRetriesExceeded
	}

	// 5. Cancelar reintentos automáticos pendientes
	for _, l := range logs {
		if l.Status == constants.SyncStatusFailed && l.NextRetryAt != nil {
			l.Status = constants.SyncStatusCancelled
			l.NextRetryAt = nil
			if err := uc.repo.UpdateInvoiceSyncLog(ctx, l); err != nil {
				uc.log.Warn(ctx).Err(err).Uint("sync_log_id", l.ID).Msg("Failed to cancel sync log before retry")
			}
		}
	}

	// 6. Crear NUEVO sync log con retry_count incrementado
	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoiceID,
		OperationType: constants.OperationTypeCreate,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    constants.MaxRetries,
		RetryCount:    lastLog.RetryCount + 1,
		TriggeredBy:   constants.TriggerManual,
	}

	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create retry sync log")
	}

	// 7. Obtener datos de la orden para construir el request
	order, err := uc.repo.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get order for retry")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("failed to get order: %w", err))
		return fmt.Errorf("failed to get order: %w", err)
	}

	// 8. Obtener configuración de facturación
	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoicing config for retry")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured)
		return errors.ErrProviderNotConfigured
	}

	// 9. Determinar integración de facturación
	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured)
		return errors.ErrProviderNotConfigured
	}

	// 10. Determinar proveedor dinámicamente según tipo de integración
	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Error al resolver proveedor de facturación para retry, usando softpymes por defecto")
		provider = dtos.ProviderSoftpymes
	}

	// 11. Construir invoiceData tipado con items de la orden
	invoiceItemDTOs := make([]dtos.InvoiceItemData, 0, len(order.Items))
	for _, item := range order.Items {
		invoiceItemDTOs = append(invoiceItemDTOs, dtos.InvoiceItemData{
			ProductID:   item.ProductID,
			SKU:         item.SKU,
			Name:        item.Name,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			Tax:         item.Tax,
			TaxRate:     item.TaxRate,
			Discount:    item.Discount,
		})
	}

	// Config específico de facturación
	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		invoiceConfigData = config.InvoiceConfig
	}

	// Inyectar URL dinámica para que el consumer seleccione entre producción y testing
	invoiceConfigData["is_testing"] = config.IsTesting
	invoiceConfigData["base_url"] = config.BaseURL
	invoiceConfigData["base_url_test"] = config.BaseURLTest

	invoiceData := dtos.InvoiceData{
		IntegrationID: integrationID,
		Customer: dtos.InvoiceCustomerData{
			Name:  invoice.CustomerName,
			Email: invoice.CustomerEmail,
			Phone: invoice.CustomerPhone,
			DNI:   invoice.CustomerDNI,
		},
		Items:        invoiceItemDTOs,
		Total:        invoice.TotalAmount,
		Subtotal:     invoice.Subtotal,
		Tax:          invoice.Tax,
		Discount:     invoice.Discount,
		ShippingCost: invoice.ShippingCost,
		Currency:     invoice.Currency,
		OrderID:      invoice.OrderID,
		Config:       invoiceConfigData,
	}

	// 12. Generar correlation ID
	correlationID := uuid.New().String()

	// 13. Construir mensaje de retry request tipado
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID:     invoice.ID,
		Provider:      provider,
		Operation:     dtos.OperationRetry,
		InvoiceData:   invoiceData,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	// 14. Publicar retry request a RabbitMQ (async)
	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("invoice_id", invoice.ID).
			Str("provider", provider).
			Msg("Failed to publish retry request to queue")

		// Marcar syncLog como failed
		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish retry to queue: " + err.Error()
		syncLog.ErrorMessage = &errorMsg

		if updateErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); updateErr != nil {
			uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
		}

		// Revertir invoice a failed
		invoice.Status = constants.InvoiceStatusFailed
		if updateErr := uc.repo.UpdateInvoice(ctx, invoice); updateErr != nil {
			uc.log.Error(ctx).Err(updateErr).Msg("Failed to revert invoice status")
		}

		return fmt.Errorf("failed to publish retry request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Int("retry_count", syncLog.RetryCount).
		Msg("Retry request published - waiting for provider response")

	// 15. Retornar éxito (invoice queda en pending, consumer lo actualizará)
	return nil
}
