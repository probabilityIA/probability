package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
)

// CreateInvoice crea una factura electrónica para una orden
func (uc *useCase) CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error) {
	uc.log.Info(ctx).Str("order_id", dto.OrderID).Msg("Creating invoice for order")

	// 1. Obtener datos de la orden
	order, err := uc.orderRepo.GetByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get order")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 2. Validar que la orden sea facturable
	if !order.Invoiceable {
		uc.log.Warn(ctx).Msg("Order is not invoiceable")
		return nil, errors.ErrOrderNotInvoiceable
	}

	// 3. Obtener configuración de facturación para la integración
	config, err := uc.configRepo.GetByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoicing config")
		return nil, errors.ErrProviderNotConfigured
	}

	if !config.Enabled {
		uc.log.Warn(ctx).Msg("Invoicing config is not enabled")
		return nil, errors.ErrConfigNotEnabled
	}

	// 4. Determinar proveedor de facturación
	var providerID uint
	if dto.InvoicingProviderID != nil {
		providerID = *dto.InvoicingProviderID
	} else {
		providerID = config.InvoicingProviderID
	}

	// 5. Verificar si ya existe una factura para esta orden y proveedor
	exists, err := uc.invoiceRepo.ExistsForOrder(ctx, order.ID, providerID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to check if invoice exists")
		return nil, fmt.Errorf("failed to check invoice existence: %w", err)
	}
	if exists {
		uc.log.Warn(ctx).Msg("Invoice already exists for order")
		return nil, errors.ErrOrderAlreadyInvoiced
	}

	// 6. Validar filtros de configuración
	if err := uc.validateInvoicingFilters(order, config); err != nil {
		uc.log.Warn(ctx).Msg("Order does not meet invoicing criteria")
		return nil, err
	}

	// 7. Obtener proveedor de facturación
	provider, err := uc.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get invoicing provider")
		return nil, errors.ErrProviderNotFound
	}

	if !provider.IsActive {
		uc.log.Warn(ctx).Msg("Provider is not active")
		return nil, errors.ErrProviderNotActive
	}

	// 8. Desencriptar credenciales
	credentials, err := uc.encryption.Decrypt(provider.Credentials)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to decrypt credentials")
		return nil, errors.ErrDecryptionFailed
	}

	// 9. Crear entidad de factura
	invoice := &entities.Invoice{
		OrderID:             order.ID,
		BusinessID:          order.BusinessID,
		InvoicingProviderID: providerID,
		Subtotal:            order.Subtotal,
		Tax:                 order.Tax,
		Discount:            order.Discount,
		ShippingCost:        order.ShippingCost,
		TotalAmount:         order.TotalAmount,
		Currency:            order.Currency,
		CustomerName:        order.CustomerName,
		CustomerEmail:       order.CustomerEmail,
		CustomerPhone:       order.CustomerPhone,
		CustomerDNI:         order.CustomerDNI,
		Status:              constants.InvoiceStatusPending,
		Notes:               dto.Notes,
		Metadata:            make(map[string]interface{}),
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
	if err := uc.invoiceRepo.Create(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice in database")
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// 12. Guardar items de factura
	for _, item := range invoiceItems {
		item.InvoiceID = invoice.ID
		if err := uc.invoiceItemRepo.Create(ctx, item); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to create invoice item")
			// No retornamos error aquí, intentamos continuar
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

	if err := uc.syncLogRepo.Create(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create sync log")
		// Continuamos aunque falle el log
	}

	// 14. Autenticar con el proveedor
	token, err := uc.providerClient.Authenticate(ctx, credentials)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to authenticate with provider")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err)
		return nil, errors.ErrAuthenticationFailed
	}

	// 15. Preparar request para el proveedor
	providerRequest := &ports.InvoiceRequest{
		Invoice:      invoice,
		InvoiceItems: invoiceItems,
		Provider:     provider,
		Config:       config.InvoiceConfig,
	}

	// 16. Enviar factura al proveedor
	response, err := uc.providerClient.CreateInvoice(ctx, token, providerRequest)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice with provider")
		uc.handleInvoiceCreationError(ctx, invoice, syncLog, err)
		return nil, errors.ErrProviderAPIError
	}

	// 17. Actualizar factura con datos del proveedor
	invoice.InvoiceNumber = response.InvoiceNumber
	invoice.ExternalID = &response.ExternalID
	invoice.InvoiceURL = response.InvoiceURL
	invoice.PDFURL = response.PDFURL
	invoice.XMLURL = response.XMLURL
	invoice.CUFE = response.CUFE
	invoice.Status = constants.InvoiceStatusIssued
	invoice.ProviderResponse = response.RawResponse

	// Parsear IssuedAt
	if response.IssuedAt != "" {
		issuedAt, err := time.Parse(time.RFC3339, response.IssuedAt)
		if err == nil {
			invoice.IssuedAt = &issuedAt
		}
	}

	// 18. Actualizar factura en BD
	if err := uc.invoiceRepo.Update(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update invoice")
		// No retornamos error, la factura ya fue creada exitosamente
	}

	// 19. Actualizar log de sincronización como exitoso
	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusSuccess
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	syncLog.ResponseStatus = 200
	syncLog.ResponseBody = response.RawResponse

	if err := uc.syncLogRepo.Update(ctx, syncLog); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update sync log")
	}

	// 20. Actualizar información de factura en la orden
	invoiceURL := ""
	if invoice.InvoiceURL != nil {
		invoiceURL = *invoice.InvoiceURL
	}
	if err := uc.orderRepo.UpdateInvoiceInfo(ctx, order.ID, invoice.InvoiceNumber, invoiceURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update order invoice info")
	}

	// 21. Publicar evento de factura creada
	if err := uc.eventPublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to publish invoice created event")
	}

	uc.log.Info(ctx).Uint("invoice_id", invoice.ID).Str("invoice_number", invoice.InvoiceNumber).Msg("Invoice created successfully")
	return invoice, nil
}

// validateInvoicingFilters valida que la orden cumpla con los filtros de configuración
func (uc *useCase) validateInvoicingFilters(order *ports.OrderData, config *entities.InvoicingConfig) error {
	// Parsear filtros
	var filters dtos.FilterConfig
	if config.Filters != nil {
		// Aquí normalmente usarías un mapper, simplificado por ahora
		if minAmount, ok := config.Filters["min_amount"].(float64); ok {
			filters.MinAmount = &minAmount
		}
		if paymentStatus, ok := config.Filters["payment_status"].(string); ok {
			filters.PaymentStatus = &paymentStatus
		}
	}

	// Validar monto mínimo
	if filters.MinAmount != nil && order.TotalAmount < *filters.MinAmount {
		return errors.ErrOrderBelowMinAmount
	}

	// Validar estado de pago
	if filters.PaymentStatus != nil && *filters.PaymentStatus == "paid" && !order.IsPaid {
		return errors.ErrOrderNotPaid
	}

	// Validar métodos de pago permitidos (si están configurados)
	if len(filters.PaymentMethods) > 0 {
		allowed := false
		for _, methodID := range filters.PaymentMethods {
			if methodID == order.PaymentMethodID {
				allowed = true
				break
			}
		}
		if !allowed {
			return errors.ErrPaymentMethodNotAllowed
		}
	}

	return nil
}

// handleInvoiceCreationError maneja errores durante la creación de factura
func (uc *useCase) handleInvoiceCreationError(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, err error) {

	// Actualizar estado de factura a failed
	invoice.Status = constants.InvoiceStatusFailed
	if updateErr := uc.invoiceRepo.Update(ctx, invoice); updateErr != nil {
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

	// Programar reintento si no se excedió el máximo
	if syncLog.RetryCount < syncLog.MaxRetries {
		nextRetry := time.Now().Add(time.Duration(constants.DefaultRetryIntervalMin) * time.Minute)
		syncLog.NextRetryAt = &nextRetry
	}

	if updateErr := uc.syncLogRepo.Update(ctx, syncLog); updateErr != nil {
		uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
	}

	// Publicar evento de factura fallida
	if publishErr := uc.eventPublisher.PublishInvoiceFailed(ctx, invoice, err.Error()); publishErr != nil {
		uc.log.Error(ctx).Err(publishErr).Msg("Failed to publish invoice failed event")
	}
}
