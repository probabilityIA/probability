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

func (uc *useCase) CreateInvoice(ctx context.Context, dto *dtos.CreateInvoiceDTO) (*entities.Invoice, error) {
	order, err := uc.repo.GetOrderByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener orden")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if !order.Invoiceable {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Orden no es facturable")
		return nil, errors.ErrOrderNotInvoiceable
	}

	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener configuración de facturación")
		return nil, errors.ErrProviderNotConfigured
	}
	if config == nil {
		config, err = uc.repo.GetEnabledConfigByBusiness(ctx, order.BusinessID)
		if err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", order.BusinessID).Msg("Error al obtener configuración de facturación por negocio")
			return nil, errors.ErrProviderNotConfigured
		}
		if config != nil && len(config.IntegrationIDs) > 0 {
			allowed := false
			for _, id := range config.IntegrationIDs {
				if id == order.IntegrationID {
					allowed = true
					break
				}
			}
			if !allowed {
				uc.log.Info(ctx).
					Str("order_id", order.ID).
					Uint("order_integration_id", order.IntegrationID).
					Uints("config_integration_ids", config.IntegrationIDs).
					Msg("Integración de origen no está en las fuentes configuradas — se omite")
				return nil, errors.ErrProviderNotConfigured
			}
		}
	}
	if config == nil {
		uc.log.Info(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Uint("integration_id", order.IntegrationID).
			Uint("business_id", order.BusinessID).
			Msg("Negocio sin configuración de facturación activa — se omite")
		return nil, errors.ErrProviderNotConfigured
	}

	if !config.Enabled {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Configuración de facturación deshabilitada")
		return nil, errors.ErrConfigNotEnabled
	}

	var integrationID uint
	if dto.InvoicingProviderID != nil {
		integrationID = *dto.InvoicingProviderID
	} else if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.log.Error(ctx).Msg("No invoicing integration configured")
		return nil, errors.ErrProviderNotConfigured
	}

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

	if err := uc.validateInvoicingFilters(order, config); err != nil {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Err(err).
			Msg("Orden no cumple criterios de facturación")
		return nil, err
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Error al resolver proveedor de facturación")
		provider = dtos.ProviderSoftpymes
	}

	invoice := &entities.Invoice{
		OrderID:                order.ID,
		BusinessID:             order.BusinessID,
		InvoicingProviderID:    nil,
		InvoicingIntegrationID: &integrationID,
		Subtotal:               order.Subtotal,
		Tax:                    order.Tax,
		Discount:               order.Discount,
		ShippingCost:           dtos.InvoiceShippingCost(order),
		ShippingDiscount:       order.ShippingDiscount,
		FreeShipping:           order.FreeShipping,
		TotalAmount:            dtos.InvoiceTotalAmount(order),
		Currency:               order.Currency,
		CustomerName:           order.CustomerName,
		CustomerEmail:          order.CustomerEmail,
		CustomerPhone:          order.CustomerPhone,
		CustomerDNI:            order.CustomerDNI,
		Status:                 constants.InvoiceStatusPending,
		IsTest:                 order.IsTest,
		Notes:                  dto.Notes,
		Metadata:               make(map[string]interface{}),
	}

	if len(order.Items) == 0 {
		uc.log.Warn(ctx).
			Str("order_id", order.ID).
			Str("order_number", order.OrderNumber).
			Msg("Orden sin items en order_items — no se puede facturar")
		return nil, fmt.Errorf("la orden %s no tiene items (order_items vacío)", order.OrderNumber)
	}

	invoiceItems := make([]*entities.InvoiceItem, 0, len(order.Items))
	for _, orderItem := range order.Items {
		unitPrice := orderItem.UnitPrice
		totalPrice := orderItem.TotalPrice
		tax := orderItem.Tax
		discount := orderItem.Discount

		if orderItem.UnitPricePresentment > 0 {
			unitPrice = orderItem.UnitPricePresentment
			totalPrice = orderItem.TotalPricePresentment
			tax = orderItem.TaxPresentment
			discount = orderItem.DiscountPresentment
		}

		unitPriceBase := orderItem.UnitPriceBase
		if orderItem.UnitPriceBasePresentment > 0 {
			unitPriceBase = orderItem.UnitPriceBasePresentment
		}

		item := &entities.InvoiceItem{
			ProductID:       orderItem.ProductID,
			SKU:             orderItem.SKU,
			Name:            orderItem.Name,
			Description:     orderItem.Description,
			Quantity:        orderItem.Quantity,
			UnitPrice:       unitPrice,
			UnitPriceBase:   unitPriceBase,
			TotalPrice:      totalPrice,
			Currency:        order.Currency,
			Tax:             tax,
			TaxRate:         orderItem.TaxRate,
			Discount:        discount,
			DiscountPercent: orderItem.DiscountPercent,
			Metadata:        make(map[string]interface{}),
		}
		invoiceItems = append(invoiceItems, item)
	}

	if err := uc.repo.CreateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create invoice in database")
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	for _, item := range invoiceItems {
		item.InvoiceID = invoice.ID
		if err := uc.repo.CreateInvoiceItem(ctx, item); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to create invoice item - cleaning up invoice")
			if delErr := uc.repo.DeleteInvoice(ctx, invoice.ID); delErr != nil {
				uc.log.Error(ctx).Err(delErr).Msg("Failed to cleanup invoice after item creation failure")
			}
			return nil, fmt.Errorf("failed to create invoice items: %w", err)
		}
	}

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
	}

	invoiceItemDTOs := make([]dtos.InvoiceItemData, 0, len(invoiceItems))
	for i, item := range invoiceItems {
		itemDTO := dtos.InvoiceItemData{
			ProductID:       item.ProductID,
			SKU:             item.SKU,
			Name:            item.Name,
			Description:     item.Description,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			UnitPriceBase:   item.UnitPriceBase,
			TotalPrice:      item.TotalPrice,
			Tax:             item.Tax,
			TaxRate:         item.TaxRate,
			Discount:        item.Discount,
			DiscountPercent: item.DiscountPercent,
		}
		if i < len(order.Items) {
			oi := order.Items[i]
			itemDTO.UnitPricePresentment = oi.UnitPricePresentment
			itemDTO.UnitPriceBasePresentment = oi.UnitPriceBasePresentment
			itemDTO.TotalPricePresentment = oi.TotalPricePresentment
			itemDTO.DiscountPresentment = oi.DiscountPresentment
			itemDTO.TaxPresentment = oi.TaxPresentment
		}
		invoiceItemDTOs = append(invoiceItemDTOs, itemDTO)
	}

	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			invoiceConfigData[k] = v
		}
	}
	invoiceConfigData["is_cod"] = order.IsCOD

	shippingTaxRate := 0.19
	if len(order.Items) > 0 && order.Items[0].TaxRate != nil && *order.Items[0].TaxRate > 0 {
		shippingTaxRate = *order.Items[0].TaxRate
	}
	shippingCostBase := dtos.ShippingCostBase(invoice.ShippingCost, shippingTaxRate)

	invoiceData := dtos.InvoiceData{
		IntegrationID: integrationID,
		Customer: dtos.InvoiceCustomerData{
			Name:  invoice.CustomerName,
			Email: invoice.CustomerEmail,
			Phone: invoice.CustomerPhone,
			DNI:   invoice.CustomerDNI,
		},
		Items:            invoiceItemDTOs,
		Total:            invoice.TotalAmount,
		Subtotal:         invoice.Subtotal,
		Tax:              invoice.Tax,
		Discount:         invoice.Discount,
		ShippingCost:     invoice.ShippingCost,
		ShippingDiscount: invoice.ShippingDiscount,
		FreeShipping:     invoice.FreeShipping,
		ShippingCostBase: shippingCostBase,
		Currency:         invoice.Currency,
		OrderID:          invoice.OrderID,
		OrderNumber:      order.OrderNumber,
		Config:           invoiceConfigData,
	}

	correlationID := uuid.New().String()

	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID:     invoice.ID,
		Provider:      provider,
		Operation:     dtos.OperationCreate,
		InvoiceData:   invoiceData,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("invoice_id", invoice.ID).
			Str("provider", provider).
			Msg("Failed to publish invoice request to queue")

		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish to queue: " + err.Error()
		syncLog.ErrorMessage = &errorMsg
		nextRetry := failedAt.Add(15 * time.Minute)
		syncLog.NextRetryAt = &nextRetry

		if updateErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); updateErr != nil {
			uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
		}

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

	if provider == dtos.ProviderSiigo {
		if enableJournal, ok := invoiceConfigData["enable_journal"].(bool); ok && enableJournal {
			journalDTO := &dtos.CreateJournalDTO{OrderID: dto.OrderID}
			if _, journalErr := uc.CreateJournal(ctx, journalDTO); journalErr != nil {
				uc.log.Warn(ctx).Err(journalErr).Msg("Auto journal creation failed (non-blocking)")
			}
		}
	}

	return invoice, nil
}

func (uc *useCase) validateInvoicingFilters(order *dtos.OrderData, config *entities.InvoicingConfig) error {
	ctx := context.Background()

	filterConfig, err := uc.parseFilterConfig(config.Filters)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to parse filter config")
		return errors.ErrInvalidFilterConfig
	}

	validators := CreateValidators(filterConfig)

	for _, validator := range validators {
		if err := validator.Validate(order); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Order failed filter validation")
			return err
		}
	}

	return nil
}

func (uc *useCase) parseFilterConfig(filtersMap map[string]interface{}) (*entities.FilterConfig, error) {
	if filtersMap == nil {
		return &entities.FilterConfig{}, nil
	}

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

func (uc *useCase) resolveProvider(ctx context.Context, integrationID uint) (string, error) {
	typeID, err := uc.repo.GetIntegrationTypeByIntegrationID(ctx, integrationID)
	if err != nil {
		return dtos.ProviderSoftpymes, err
	}

	switch typeID {
	case 5:
		return dtos.ProviderSoftpymes, nil
	case 7:
		return dtos.ProviderFactus, nil
	case 8:
		return dtos.ProviderSiigo, nil
	default:
		uc.log.Warn(ctx).
			Uint("integration_id", integrationID).
			Int("type_id", typeID).
			Msg("Unknown integration type for invoicing, defaulting to softpymes")
		return dtos.ProviderSoftpymes, nil
	}
}

func (uc *useCase) handleInvoiceCreationError(ctx context.Context, invoice *entities.Invoice, syncLog *entities.InvoiceSyncLog, err error) {

	invoice.Status = constants.InvoiceStatusFailed
	if updateErr := uc.repo.UpdateInvoice(ctx, invoice); updateErr != nil {
		uc.log.Error(ctx).Err(updateErr).Msg("Failed to update invoice status to failed")
	}

	completedAt := time.Now()
	duration := int(completedAt.Sub(syncLog.StartedAt).Milliseconds())
	syncLog.Status = constants.SyncStatusFailed
	syncLog.CompletedAt = &completedAt
	syncLog.Duration = &duration
	errorMsg := err.Error()
	syncLog.ErrorMessage = &errorMsg

	if syncLog.RetryCount < syncLog.MaxRetries && syncLog.Status != constants.SyncStatusCancelled {
		nextRetry := time.Now().Add(time.Duration(constants.DefaultRetryIntervalMin) * time.Minute)
		syncLog.NextRetryAt = &nextRetry
	}

	if updateErr := uc.repo.UpdateInvoiceSyncLog(ctx, syncLog); updateErr != nil {
		uc.log.Error(ctx).Err(updateErr).Msg("Failed to update sync log")
	}

	if publishErr := uc.eventPublisher.PublishInvoiceFailed(ctx, invoice, err.Error()); publishErr != nil {
		uc.log.Error(ctx).Err(publishErr).Msg("Failed to publish invoice failed event")
	}

	if publishErr := uc.ssePublisher.PublishInvoiceFailed(ctx, invoice, err.Error()); publishErr != nil {
		uc.log.Error(ctx).Err(publishErr).Msg("Failed to publish invoice failed SSE event")
	}
}

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
