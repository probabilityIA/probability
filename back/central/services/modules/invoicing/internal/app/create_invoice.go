package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// CreateInvoice crea una factura electrónica para una orden
func (uc *useCase) CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error) {
	// 1. Obtener datos de la orden
	order, err := uc.repo.GetOrderByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener orden")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 2. Validar que la orden sea facturable
	if !order.Invoiceable {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Orden no es facturable")
		return nil, errors.ErrOrderNotInvoiceable
	}

	// 3. Obtener configuración de facturación para la integración
	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener configuración de facturación")
		return nil, errors.ErrProviderNotConfigured
	}

	if !config.Enabled {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Configuración de facturación deshabilitada")
		return nil, errors.ErrConfigNotEnabled
	}

	// Validar auto_invoice solo para facturas automáticas
	if !dto.IsManual && !config.AutoInvoice {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Facturación automática deshabilitada")
		return nil, errors.ErrAutoInvoiceNotEnabled
	}

	// 4. Determinar integración de facturación
	var integrationID uint
	if dto.InvoicingProviderID != nil {
		// Dual-read: Si se proporciona el ID viejo, usarlo temporalmente
		integrationID = *dto.InvoicingProviderID
	} else if config.InvoicingIntegrationID != nil {
		// Usar el nuevo campo de integración
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		// Fallback al campo viejo durante migración
		integrationID = *config.InvoicingProviderID
	} else {
		uc.log.Error(ctx).Msg("No invoicing integration configured")
		return nil, errors.ErrProviderNotConfigured
	}

	// 5. Verificar si ya existe una factura para esta orden e integración
	exists, err := uc.repo.InvoiceExistsForOrder(ctx, order.ID, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al verificar factura existente")
		return nil, fmt.Errorf("failed to check invoice existence: %w", err)
	}
	if exists {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Ya existe factura para esta orden")
		return nil, errors.ErrOrderAlreadyInvoiced
	}

	// 6. Validar filtros de configuración
	if err := uc.validateInvoicingFilters(order, config); err != nil {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Err(err).
			Msg("Orden no cumple criterios de facturación")
		return nil, err
	}

	// 7. Determinar proveedor de facturación
	provider := dtos.ProviderSoftpymes

	// 8. Crear entidad de factura (estado pending)
	invoice := &entities.Invoice{
		OrderID:                order.ID,
		BusinessID:             order.BusinessID,
		InvoicingProviderID:    nil,
		InvoicingIntegrationID: &integrationID,
		Subtotal:               order.Subtotal,
		Tax:                    order.Tax,
		Discount:               order.Discount,
		ShippingCost:           order.ShippingCost,
		TotalAmount:            order.TotalAmount,
		Currency:               order.Currency,
		CustomerName:           order.CustomerName,
		CustomerEmail:          order.CustomerEmail,
		CustomerPhone:          order.CustomerPhone,
		CustomerDNI:            order.CustomerDNI,
		Status:                 constants.InvoiceStatusPending,
		Notes:                  dto.Notes,
		Metadata:               make(map[string]interface{}),
	}

	// 9. Crear items de factura desde los items de la orden
	invoiceItems := make([]*entities.InvoiceItem, 0, len(order.Items))
	for _, orderItem := range order.Items {
		item := &entities.InvoiceItem{
			ProductID:   orderItem.ProductID,
			SKU:         orderItem.SKU,
			Name:        orderItem.Name,
			Description: orderItem.Description,
			Quantity:    orderItem.Quantity,
			UnitPrice:   orderItem.UnitPrice,
			TotalPrice:  orderItem.TotalPrice,
			Currency:    order.Currency,
			Tax:         orderItem.Tax,
			TaxRate:     orderItem.TaxRate,
			Discount:    orderItem.Discount,
			Metadata:    make(map[string]interface{}),
		}
		invoiceItems = append(invoiceItems, item)
	}

	// 10. Guardar factura en BD (estado pending)
	if err := uc.repo.CreateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice in database")
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// 11. Guardar items de factura
	for _, item := range invoiceItems {
		item.InvoiceID = invoice.ID
		if err := uc.repo.CreateInvoiceItem(ctx, item); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to create invoice item - cleaning up invoice")
			// Cleanup: eliminar la factura incompleta para evitar items parciales
			if delErr := uc.repo.DeleteInvoice(ctx, invoice.ID); delErr != nil {
				uc.log.Error(ctx).Err(delErr).Msg("Failed to cleanup invoice after item creation failure")
			}
			return nil, fmt.Errorf("failed to create invoice items: %w", err)
		}
	}

	// 12. Crear log de sincronización
	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeCreate,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    constants.MaxRetries,
		RetryCount:    0,
	}

	if dto.IsManual {
		syncLog.TriggeredBy = constants.TriggerManual
		syncLog.UserID = dto.CreatedByUserID
	} else {
		syncLog.TriggeredBy = constants.TriggerAuto
	}

	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create sync log")
		// Continuamos aunque falle el log
	}

	// 13. Preparar datos de facturación tipados para el proveedor
	invoiceItemDTOs := make([]dtos.InvoiceItemData, 0, len(invoiceItems))
	for _, item := range invoiceItems {
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

	// Config específico de facturación (invoice_config desde DB)
	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		invoiceConfigData = config.InvoiceConfig
	}

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

	// 14. Generar correlation ID único para request/response
	correlationID := uuid.New().String()

	// 15. Construir mensaje de request tipado
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID:     invoice.ID,
		Provider:      provider,
		Operation:     dtos.OperationCreate,
		InvoiceData:   invoiceData,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	// 16. Publicar request a RabbitMQ (async)
	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("invoice_id", invoice.ID).
			Str("provider", provider).
			Msg("Failed to publish invoice request to queue")

		// Marcar syncLog como failed
		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish to queue: " + err.Error()
		syncLog.ErrorMessage = &errorMsg

		if updateErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); updateErr != nil {
			uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
		}

		// Marcar invoice como failed
		invoice.Status = constants.InvoiceStatusFailed
		if updateErr := uc.repo.UpdateInvoice(ctx, invoice); updateErr != nil {
			uc.log.Error(ctx).Err(updateErr).Msg("Failed to update invoice status")
		}

		return nil, fmt.Errorf("failed to publish invoice request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("Invoice request published - waiting for provider response")

	// 17. Retornar invoice inmediatamente (estado pending)
	return invoice, nil
}

// validateInvoicingFilters valida que la orden cumpla con los filtros de configuración
func (uc *useCase) validateInvoicingFilters(order *dtos.OrderData, config *entities.InvoicingConfig) error {
	ctx := context.Background()

	// 1. Parsear configuración de filtros desde JSON
	filterConfig, err := uc.parseFilterConfig(config.Filters)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to parse filter config")
		return errors.ErrInvalidFilterConfig
	}

	// 2. Crear validadores dinámicamente
	validators := CreateValidators(filterConfig)

	// 3. Ejecutar todas las validaciones
	for _, validator := range validators {
		if err := validator.Validate(order); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Order failed filter validation")
			return err
		}
	}

	return nil
}

// parseFilterConfig convierte el map[string]interface{} a entities.FilterConfig estructurado
func (uc *useCase) parseFilterConfig(filtersMap map[string]interface{}) (*entities.FilterConfig, error) {
	if filtersMap == nil {
		return &entities.FilterConfig{}, nil
	}

	// Usar JSON marshal/unmarshal para conversión segura
	jsonData, err := json.Marshal(filtersMap)
	if err != nil {
		return nil, err
	}

	var config entities.FilterConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// handleInvoiceCreationError maneja errores durante la creación de factura
func (uc *useCase) handleInvoiceCreationError(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, err error) {

	// Actualizar estado de factura a failed
	invoice.Status = constants.InvoiceStatusFailed
	if updateErr := uc.repo.UpdateInvoice(ctx, invoice); updateErr != nil {
		uc.log.Error(ctx).Err(updateErr).Msg("Failed to update invoice status to failed")
	}

	// Actualizar sync log
	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusFailed
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	errorMsg := err.Error()
	syncLog.ErrorMessage = &errorMsg

	// Programar reintento si no se excedió el máximo y no fue cancelado previamente
	if syncLog.RetryCount < syncLog.MaxRetries && syncLog.Status != constants.SyncStatusCancelled {
		nextRetry := time.Now().Add(time.Duration(constants.DefaultRetryIntervalMin) * time.Minute)
		syncLog.NextRetryAt = &nextRetry
	}

	if updateErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); updateErr != nil {
		uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
	}

	// Publicar evento de factura fallida (RabbitMQ)
	if publishErr := uc.eventPublisher.PublishInvoiceFailed(ctx, invoice, err.Error()); publishErr != nil {
		uc.log.Error(ctx).Err(publishErr).Msg("Failed to publish invoice failed event")
	}

	// Publicar evento SSE en tiempo real (Redis Pub/Sub)
	if publishErr := uc.ssePublisher.PublishInvoiceFailed(ctx, invoice, err.Error()); publishErr != nil {
		uc.log.Error(ctx).Err(publishErr).Msg("Failed to publish invoice failed SSE event")
	}
}

// updateInvoiceWithRetry reintenta UpdateInvoice hasta 3 veces tras éxito del proveedor.
// Si falla todas las veces, guarda fallbackData en el sync log como respaldo
// para recuperación manual, evitando facturas fantasma.
func (uc *useCase) updateInvoiceWithRetry(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, fallbackData map[string]interface{}) error {
	var updateErr error
	for attempt := 0; attempt < 3; attempt++ {
		updateErr = uc.repo.UpdateInvoice(ctx, invoice)
		if updateErr == nil {
			return nil
		}
		uc.log.Warn(ctx).Err(updateErr).Int("attempt", attempt+1).Msg("UpdateInvoice failed, retrying...")
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
	}

	// Todas las reintentos fallaron - guardar datos del proveedor en sync log como fallback
	uc.log.Error(ctx).Err(updateErr).Msg("CRITICAL: Invoice created at provider but DB update failed after 3 attempts")
	criticalMsg := "CRITICAL: Invoice created at provider but DB update failed"
	syncLog.ErrorMessage = &criticalMsg
	syncLog.ResponseBody = fallbackData
	syncLog.Status = constants.SyncStatusFailed
	if logErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); logErr != nil {
		uc.log.Error(ctx).Err(logErr).Msg("Failed to save fallback data in sync log")
	}
	return fmt.Errorf("invoice created at provider but failed to save: %w", updateErr)
}
