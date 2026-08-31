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

func (uc *useCase) CreateJournal(ctx context.Context, dto *dtos.CreateJournalDTO) (*entities.Invoice, error) {
	order, err := uc.repo.GetOrderByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener orden para journal")
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al obtener configuración de facturación para journal")
		return nil, errors.ErrProviderNotConfigured
	}
	if config == nil {
		config, err = uc.repo.GetEnabledConfigByBusiness(ctx, order.BusinessID)
		if err != nil {
			uc.log.Error(ctx).Err(err).Uint("business_id", order.BusinessID).Msg("Error al obtener config por negocio para journal")
			return nil, errors.ErrProviderNotConfigured
		}
	}
	if config == nil || !config.Enabled {
		uc.log.Info(ctx).Str("order_id", order.ID).Msg("Sin configuración activa para journal")
		return nil, errors.ErrProviderNotConfigured
	}

	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		return nil, errors.ErrProviderNotConfigured
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al resolver proveedor para journal")
		return nil, errors.ErrProviderNotConfigured
	}
	if provider != dtos.ProviderSiigo {
		uc.log.Info(ctx).Str("provider", provider).Msg("Journal solo aplica para Siigo, omitiendo")
		return nil, fmt.Errorf("journals are only supported for Siigo provider, got: %s", provider)
	}

	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		invoiceConfigData = config.InvoiceConfig
	}
	enableJournal, _ := invoiceConfigData["enable_journal"].(bool)
	inventoryExitOnly, _ := invoiceConfigData["inventory_exit_only"].(bool)
	if !enableJournal && !inventoryExitOnly {
		uc.log.Info(ctx).Str("order_id", order.ID).Msg("Comprobante contable no habilitado en config, omitiendo")
		return nil, fmt.Errorf("journal is not enabled in invoicing config")
	}
	if inventoryExitOnly {
		inventoryAccount, _ := invoiceConfigData["inventory_exit_account_code"].(string)
		offsetAccount, _ := invoiceConfigData["inventory_exit_offset_account_code"].(string)
		if inventoryAccount == "" || offsetAccount == "" {
			uc.log.Warn(ctx).Str("order_id", order.ID).Msg("Salida de inventario sin cuentas contables configuradas")
			return nil, fmt.Errorf("la salida de inventario requiere inventory_exit_account_code y inventory_exit_offset_account_code")
		}
	}

	if len(order.Items) == 0 {
		uc.log.Warn(ctx).Str("order_id", order.ID).Msg("Orden sin items — no se puede crear journal")
		return nil, fmt.Errorf("la orden %s no tiene items (order_items vacío)", order.OrderNumber)
	}

	invoice := &entities.Invoice{
		OrderID:                order.ID,
		BusinessID:             order.BusinessID,
		InvoicingIntegrationID: &integrationID,
		Subtotal:               order.Subtotal,
		Tax:                    order.Tax,
		Discount:               order.Discount,
		ShippingCost:           dtos.InvoiceShippingCost(order),
		TotalAmount:            dtos.InvoiceTotalAmount(order),
		Currency:               order.Currency,
		CustomerName:           order.CustomerName,
		CustomerEmail:          order.CustomerEmail,
		CustomerPhone:          order.CustomerPhone,
		CustomerDNI:            order.CustomerDNI,
		Status:                 constants.InvoiceStatusPending,
		Metadata: map[string]interface{}{
			"type":               "journal",
			"original_operation": dtos.OperationCreateJournal,
		},
	}

	if err := uc.repo.CreateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create journal invoice record")
		return nil, fmt.Errorf("failed to create journal invoice: %w", err)
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

		item := &entities.InvoiceItem{
			InvoiceID:       invoice.ID,
			ProductID:       orderItem.ProductID,
			SKU:             orderItem.SKU,
			Name:            orderItem.Name,
			Description:     orderItem.Description,
			Quantity:        orderItem.Quantity,
			UnitPrice:       unitPrice,
			TotalPrice:      totalPrice,
			Currency:        order.Currency,
			Tax:             tax,
			TaxRate:         orderItem.TaxRate,
			Discount:        discount,
			DiscountPercent: orderItem.DiscountPercent,
			Metadata:        make(map[string]interface{}),
		}
		invoiceItems = append(invoiceItems, item)

		if err := uc.repo.CreateInvoiceItem(ctx, item); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to create journal invoice item")
			if delErr := uc.repo.DeleteInvoice(ctx, invoice.ID); delErr != nil {
				uc.log.Error(ctx).Err(delErr).Msg("Failed to cleanup journal invoice after item creation failure")
			}
			return nil, fmt.Errorf("failed to create journal invoice items: %w", err)
		}
	}

	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeCreateJournal,
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
		uc.log.Error(ctx).Err(err).Msg("Failed to create journal sync log")
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
			TotalPrice:      item.TotalPrice,
			Tax:             item.Tax,
			TaxRate:         item.TaxRate,
			Discount:        item.Discount,
			DiscountPercent: item.DiscountPercent,
		}
		if i < len(order.Items) {
			oi := order.Items[i]
			itemDTO.UnitPricePresentment = oi.UnitPricePresentment
			itemDTO.TotalPricePresentment = oi.TotalPricePresentment
			itemDTO.DiscountPresentment = oi.DiscountPresentment
			itemDTO.TaxPresentment = oi.TaxPresentment
		}
		invoiceItemDTOs = append(invoiceItemDTOs, itemDTO)
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
		OrderNumber:  order.OrderNumber,
		Config:       invoiceConfigData,
	}

	correlationID := uuid.New().String()
	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID:     invoice.ID,
		Provider:      provider,
		Operation:     dtos.OperationCreateJournal,
		InvoiceData:   invoiceData,
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to publish journal request to queue")

		failedAt := time.Now()
		duration := int(failedAt.Sub(syncLog.StartedAt).Milliseconds())
		syncLog.Status = constants.SyncStatusFailed
		syncLog.CompletedAt = &failedAt
		syncLog.Duration = &duration
		errorMsg := "Failed to publish to queue: " + err.Error()
		syncLog.ErrorMessage = &errorMsg
		_ = uc.repo.UpdateInvoiceSyncLog(ctx, syncLog)

		invoice.Status = constants.InvoiceStatusFailed
		_ = uc.repo.UpdateInvoice(ctx, invoice)

		return nil, fmt.Errorf("failed to publish journal request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("correlation_id", correlationID).
		Msg("Journal request published - waiting for Siigo response")

	return invoice, nil
}
