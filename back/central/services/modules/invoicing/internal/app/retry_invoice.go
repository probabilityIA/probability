package app

import (
	"context"
	"fmt"
<<<<<<< HEAD

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// RetryInvoice reintenta la creación de una factura fallida
func (uc *useCase) RetryInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Retrying invoice creation")

	// 1. Obtener factura
	invoice, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
=======
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	softpymesBundle "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// RetryInvoice reintenta la creación de una factura fallida (in-place, sin eliminar)
func (uc *useCase) RetryInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Retrying invoice creation")

	// 1. Obtener factura existente (NO se elimina)
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

<<<<<<< HEAD
	// 2. Validar que esté en estado failed
=======
	// 2. Validar que esté en estado failed (solo se puede reintentar desde failed)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if invoice.Status != constants.InvoiceStatusFailed {
		return errors.ErrRetryNotAllowed
	}

<<<<<<< HEAD
	// 3. Obtener logs de sincronización
	logs, err := uc.syncLogRepo.GetByInvoiceID(ctx, invoiceID)
=======
	// 2.5. Lock optimista: marcar invoice como pending ANTES de llamar al proveedor.
	// Si el RetryConsumer o un retry manual concurrente intenta procesar la misma factura,
	// verá status != failed y saldrá con ErrRetryNotAllowed (paso 2).
	invoice.Status = constants.InvoiceStatusPending
	if err := uc.repo.UpdateInvoice(ctx, invoice); err != nil {
		return fmt.Errorf("failed to lock invoice for retry: %w", err)
	}

	// 3. Obtener logs de sincronización para verificar reintentos
	logs, err := uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil || len(logs) == 0 {
		return fmt.Errorf("no sync logs found for invoice")
	}

<<<<<<< HEAD
	lastLog := logs[len(logs)-1]
=======
	lastLog := logs[0] // Ordenados por created_at DESC
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	// 4. Validar que no se haya excedido el máximo de reintentos
	if lastLog.RetryCount >= lastLog.MaxRetries {
		return errors.ErrMaxRetriesExceeded
	}

<<<<<<< HEAD
	// 5. Reintentar creación usando CreateInvoice
	dto := &dtos.CreateInvoiceDTO{
		OrderID:  invoice.OrderID,
		Notes:    invoice.Notes,
		IsManual: true, // Los reintentos se consideran manuales
	}

	_, err = uc.CreateInvoice(ctx, dto)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Retry failed")
		return err
	}

	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Invoice retry successful")
=======
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
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("failed to get order: %w", err), nil)
		return fmt.Errorf("failed to get order: %w", err)
	}

	// 8. Obtener configuración de facturación
	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoicing config for retry")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured, nil)
		return errors.ErrProviderNotConfigured
	}

	// 9. Determinar integración de facturación
	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotConfigured, nil)
		return errors.ErrProviderNotConfigured
	}

	// 10. Obtener integración y credenciales
	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrProviderNotFound, nil)
		return errors.ErrProviderNotFound
	}

	credentialsMap := make(map[string]interface{})
	if integration.Config != nil {
		if configMap, ok := integration.Config.(map[string]interface{}); ok {
			credentialsMap = configMap
		}
	}

	apiKey, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrDecryptionFailed, nil)
		return errors.ErrDecryptionFailed
	}

	apiSecret, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, errors.ErrDecryptionFailed, nil)
		return errors.ErrDecryptionFailed
	}

	credentialsMap["api_key"] = apiKey
	credentialsMap["api_secret"] = apiSecret

	// 11. Construir invoiceData con items de la orden
	items := make([]map[string]interface{}, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, map[string]interface{}{
			"product_id":  item.ProductID,
			"sku":         item.SKU,
			"name":        item.Name,
			"description": item.Description,
			"quantity":    item.Quantity,
			"unit_price":  item.UnitPrice,
			"total_price": item.TotalPrice,
			"tax":         item.Tax,
			"tax_rate":    item.TaxRate,
			"discount":    item.Discount,
		})
	}

	invoiceData := map[string]interface{}{
		"credentials":   credentialsMap,
		"customer":      map[string]interface{}{
			"name":  invoice.CustomerName,
			"email": invoice.CustomerEmail,
			"phone": invoice.CustomerPhone,
			"dni":   invoice.CustomerDNI,
		},
		"items":         items,
		"total":         invoice.TotalAmount,
		"subtotal":      invoice.Subtotal,
		"tax":           invoice.Tax,
		"discount":      invoice.Discount,
		"shipping_cost": invoice.ShippingCost,
		"currency":      invoice.Currency,
		"order_id":      invoice.OrderID,
		"config":        integration.Config,
	}

	// 12. Obtener bundle del proveedor y llamar
	if integration.IntegrationType != integrationCore.IntegrationTypeInvoicing {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("integration type mismatch"), nil)
		return errors.ErrProviderNotFound
	}

	integrationBundle, ok := uc.integrationCore.GetRegisteredIntegration(integrationCore.IntegrationTypeInvoicing)
	if !ok {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("invoicing bundle not found"), nil)
		return errors.ErrProviderNotFound
	}

	softpymes, ok := integrationBundle.(*softpymesBundle.Bundle)
	if !ok {
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("invalid bundle type"), nil)
		return errors.ErrProviderNotFound
	}

	err = softpymes.CreateInvoice(ctx, invoiceData)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Retry failed - provider error")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err, invoiceData)
		return errors.ErrProviderAPIError
	}

	// 13. Éxito - Actualizar factura con datos del proveedor
	if invoiceNumber, ok := invoiceData["invoice_number"].(string); ok {
		invoice.InvoiceNumber = invoiceNumber
	}
	if externalID, ok := invoiceData["external_id"].(string); ok {
		invoice.ExternalID = &externalID
	}
	if invoiceURL, ok := invoiceData["invoice_url"].(string); ok {
		invoice.InvoiceURL = &invoiceURL
	}
	if pdfURL, ok := invoiceData["pdf_url"].(string); ok {
		invoice.PDFURL = &pdfURL
	}
	if xmlURL, ok := invoiceData["xml_url"].(string); ok {
		invoice.XMLURL = &xmlURL
	}
	if cufe, ok := invoiceData["cufe"].(string); ok {
		invoice.CUFE = &cufe
	}

	invoice.Status = constants.InvoiceStatusIssued

	if issuedAtStr, ok := invoiceData["issued_at"].(string); ok && issuedAtStr != "" {
		issuedAt, err := time.Parse(time.RFC3339, issuedAtStr)
		if err == nil {
			invoice.IssuedAt = &issuedAt
		}
	}

	if err := uc.updateInvoiceWithRetry(ctx, invoice, syncLog, invoiceData); err != nil {
		return err
	}

	// 14. Actualizar sync log como exitoso
	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusSuccess
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	invoice.ProviderResponse = invoiceData

	uc.populateSyncLogAudit(syncLog, invoiceData)

	if err := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update sync log after retry success")
	}

	// 15. Actualizar información de factura en la orden
	invoiceURL := ""
	if invoice.InvoiceURL != nil {
		invoiceURL = *invoice.InvoiceURL
	}
	if err := uc.repo.UpdateOrderInvoiceInfo(ctx, order.ID, invoice.InvoiceNumber, invoiceURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update order invoice info after retry")
	}

	// 16. Publicar eventos
	if err := uc.eventPublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish invoice created event after retry")
	}
	if err := uc.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish invoice created SSE event after retry")
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Int("retry_count", syncLog.RetryCount).
		Msg("Invoice retry completed successfully")

>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	return nil
}
