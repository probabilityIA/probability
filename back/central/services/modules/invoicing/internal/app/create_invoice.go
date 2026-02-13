package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	integrationCore "github.com/secamc93/probability/back/central/services/integrations/core"
	softpymesBundle "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes"
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
			Msg("❌ Orden no es facturable")
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
			Msg("❌ Configuración de facturación deshabilitada")
		return nil, errors.ErrConfigNotEnabled
	}

	// Validar auto_invoice solo para facturas automáticas
	if !dto.IsManual && !config.AutoInvoice {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("❌ Facturación automática deshabilitada")
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
			Msg("❌ Ya existe factura para esta orden")
		return nil, errors.ErrOrderAlreadyInvoiced
	}

	// 6. Validar filtros de configuración
	if err := uc.validateInvoicingFilters(order, config); err != nil {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Err(err).
			Msg("❌ Orden no cumple criterios de facturación")
		return nil, err
	}

	// 7. Obtener integración de facturación desde integrationCore
	integrationIDStr := fmt.Sprintf("%d", integrationID)
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationIDStr)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoicing integration")
		return nil, errors.ErrProviderNotFound
	}

	// 8. Desencriptar credenciales usando integrationCore
	// NOTA: Las credenciales ahora se manejan internamente por el bundle de softpymes
	// No necesitamos desencriptarlas aquí

	// 9. Crear entidad de factura
	invoice := &entities.Invoice{
		OrderID:                order.ID,
		BusinessID:             order.BusinessID,
		InvoicingProviderID:    nil,            // NULL - campo legacy deprecado (FK hacia invoicing_providers)
		InvoicingIntegrationID: &integrationID, // Campo actual (FK hacia integrations)
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

	// 10. Crear items de factura desde los items de la orden
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

	// 11. Guardar factura en BD (estado pending)
	if err := uc.repo.CreateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice in database")
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// 12. Guardar items de factura
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

	// 13. Crear log de sincronización
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

	// 14. Obtener credentials de la integración usando integrationCore
	credentialsMap := make(map[string]interface{})
	if integration.Config != nil {
		if configMap, ok := integration.Config.(map[string]interface{}); ok {
			credentialsMap = configMap
		}
	}

	// Desencriptar credenciales específicas si es necesario
	// NOTA: El bundle de Softpymes maneja la autenticación internamente
	apiKey, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_key")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt api_key")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err, nil)
		return nil, errors.ErrDecryptionFailed
	}

	apiSecret, err := uc.integrationCore.DecryptCredential(ctx, integrationIDStr, "api_secret")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt api_secret")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err, nil)
		return nil, errors.ErrDecryptionFailed
	}

	credentialsMap["api_key"] = apiKey
	credentialsMap["api_secret"] = apiSecret

	// 15. Preparar datos para el bundle de Softpymes
	invoiceItems2 := make([]map[string]interface{}, 0, len(invoiceItems))
	for _, item := range invoiceItems {
		itemMap := map[string]interface{}{
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
		}
		invoiceItems2 = append(invoiceItems2, itemMap)
	}

	customerData := map[string]interface{}{
		"name":  invoice.CustomerName,
		"email": invoice.CustomerEmail,
		"phone": invoice.CustomerPhone,
		"dni":   invoice.CustomerDNI,
	}

	invoiceData := map[string]interface{}{
		"credentials":   credentialsMap,
		"customer":      customerData,
		"items":         invoiceItems2,
		"total":         invoice.TotalAmount,
		"subtotal":      invoice.Subtotal,
		"tax":           invoice.Tax,
		"discount":      invoice.Discount,
		"shipping_cost": invoice.ShippingCost,
		"currency":      invoice.Currency,
		"order_id":      invoice.OrderID,
		"config":        integration.Config, // Config de la integración (contiene referer, api_url, etc)
	}

	// 16. Enviar factura al proveedor usando el bundle de Softpymes
	// Obtener el bundle registrado desde integrationCore
	if integration.IntegrationType != integrationCore.IntegrationTypeInvoicing {
		uc.log.Error(ctx).Msg("Integration is not an invoicing type")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("integration type mismatch"), nil)
		return nil, errors.ErrProviderNotFound
	}

	integrationBundle, ok := uc.integrationCore.GetRegisteredIntegration(integrationCore.IntegrationTypeInvoicing)
	if !ok {
		uc.log.Error(ctx).Msg("Invoicing integration bundle not registered")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("invoicing bundle not found"), nil)
		return nil, errors.ErrProviderNotFound
	}

	// Cast al bundle de Softpymes
	softpymes, ok := integrationBundle.(*softpymesBundle.Bundle)
	if !ok {
		uc.log.Error(ctx).Msg("Failed to cast integration to softpymes bundle")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, fmt.Errorf("invalid bundle type"), nil)
		return nil, errors.ErrProviderNotFound
	}

	err = softpymes.CreateInvoice(ctx, invoiceData)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice with provider")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err, invoiceData)
		return nil, errors.ErrProviderAPIError
	}

	// 17. Actualizar factura con datos del proveedor
	// El bundle de Softpymes actualiza invoiceData con la respuesta
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

	// Parsear IssuedAt si está presente
	if issuedAtStr, ok := invoiceData["issued_at"].(string); ok && issuedAtStr != "" {
		issuedAt, err := time.Parse(time.RFC3339, issuedAtStr)
		if err == nil {
			invoice.IssuedAt = &issuedAt
		}
	}

	// 18. Actualizar factura en BD (con reintentos para evitar factura fantasma)
	if err := uc.updateInvoiceWithRetry(ctx, invoice, syncLog, invoiceData); err != nil {
		return nil, err
	}

	// 19. Actualizar log de sincronización como exitoso
	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusSuccess
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	// La respuesta completa está ahora en invoiceData
	invoice.ProviderResponse = invoiceData

	// Extraer audit data del invoiceData (capturada por el cliente Softpymes)
	uc.populateSyncLogAudit(syncLog, invoiceData)

	if err := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update sync log")
	}

	// 20. Actualizar información de factura en la orden
	invoiceURL := ""
	if invoice.InvoiceURL != nil {
		invoiceURL = *invoice.InvoiceURL
	}
	if err := uc.repo.UpdateOrderInvoiceInfo(ctx, order.ID, invoice.InvoiceNumber, invoiceURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update order invoice info")
	}

	// 21. Publicar evento de factura creada (RabbitMQ)
	if err := uc.eventPublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish invoice created event")
	}

	// 22. Publicar evento SSE en tiempo real (Redis Pub/Sub)
	if err := uc.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish invoice created SSE event")
	}

	uc.log.Info(ctx).Uint("invoice_id", invoice.ID).Str("invoice_number", invoice.InvoiceNumber).Msg("Invoice created successfully")
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
func (uc *useCase) handleInvoiceCreationError(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, err error, invoiceData map[string]interface{}) {

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

	// Extraer audit data si existe
	uc.populateSyncLogAudit(syncLog, invoiceData)

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
// Si falla todas las veces, guarda los datos del proveedor en el sync log como fallback
// para recuperación manual, evitando facturas fantasma.
func (uc *useCase) updateInvoiceWithRetry(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, invoiceData map[string]interface{}) error {
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
	syncLog.ResponseBody = invoiceData
	syncLog.Status = constants.SyncStatusFailed
	if logErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); logErr != nil {
		uc.log.Error(ctx).Err(logErr).Msg("Failed to save fallback data in sync log")
	}
	return fmt.Errorf("invoice created at provider but failed to save: %w", updateErr)
}

// populateSyncLogAudit extrae audit data del invoiceData y la almacena en el sync log
func (uc *useCase) populateSyncLogAudit(syncLog *entities.InvoiceSyncLog, invoiceData map[string]interface{}) {
	if invoiceData == nil {
		return
	}
	auditData, ok := invoiceData["_audit"].(map[string]interface{})
	if !ok {
		return
	}
	if reqURL, ok := auditData["request_url"].(string); ok {
		syncLog.RequestURL = reqURL
	}
	if reqPayload, ok := auditData["request_payload"].(map[string]interface{}); ok {
		syncLog.RequestPayload = reqPayload
	}
	if respStatus, ok := auditData["response_status"].(int); ok {
		syncLog.ResponseStatus = respStatus
	}
	if respBody, ok := auditData["response_body"].(string); ok {
		var bodyMap map[string]interface{}
		if json.Unmarshal([]byte(respBody), &bodyMap) == nil {
			syncLog.ResponseBody = bodyMap
		}
	}
}
