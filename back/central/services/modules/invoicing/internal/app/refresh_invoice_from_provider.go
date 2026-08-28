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

func (uc *useCase) buildCheckStatusMessage(ctx context.Context, invoice *entities.Invoice) (*dtos.InvoiceRequestMessage, error) {
	order, err := uc.repo.GetOrderByID(ctx, invoice.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		return nil, errors.ErrProviderNotConfigured
	}
	if config == nil {
		config, err = uc.repo.GetEnabledConfigByBusiness(ctx, order.BusinessID)
		if err != nil || config == nil {
			return nil, errors.ErrProviderNotConfigured
		}
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
		provider = dtos.ProviderSoftpymes
	}

	invoiceConfigData := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		invoiceConfigData = config.InvoiceConfig
	}
	invoiceConfigData["is_testing"] = config.IsTesting
	invoiceConfigData["base_url"] = config.BaseURL
	invoiceConfigData["base_url_test"] = config.BaseURLTest

	if invoice.ExternalID != nil && *invoice.ExternalID != "" {
		invoiceConfigData["external_id"] = *invoice.ExternalID
	}

	return &dtos.InvoiceRequestMessage{
		InvoiceID: invoice.ID,
		Provider:  provider,
		Operation: dtos.OperationCheckStatus,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			OrderID:       invoice.OrderID,
			OrderNumber:   order.OrderNumber,
			Currency:      invoice.Currency,
			Config:        invoiceConfigData,
		},
		CorrelationID: uuid.New().String(),
		Timestamp:     time.Now(),
	}, nil
}

func (uc *useCase) RefreshInvoiceFromProvider(ctx context.Context, invoiceID uint) error {
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	if invoice.Status != constants.InvoiceStatusIssued && invoice.Status != constants.InvoiceStatusPending {
		return errors.ErrRetryNotAllowed
	}
	if invoice.ExternalID == nil || *invoice.ExternalID == "" {
		return fmt.Errorf("la factura no tiene identificador del proveedor: no se puede consultar")
	}

	syncLog := &entities.InvoiceSyncLog{
		InvoiceID:     invoice.ID,
		OperationType: constants.OperationTypeQuery,
		Status:        constants.SyncStatusProcessing,
		StartedAt:     time.Now(),
		MaxRetries:    constants.MaxCheckAttempts,
		TriggeredBy:   constants.TriggerManual,
	}
	if err := uc.repo.CreateInvoiceSyncLog(ctx, syncLog); err != nil {
		uc.log.Warn(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to create refresh sync log")
	}

	message, err := uc.buildCheckStatusMessage(ctx, invoice)
	if err != nil {
		return err
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, message); err != nil {
		return fmt.Errorf("failed to publish refresh request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("provider", message.Provider).
		Str("correlation_id", message.CorrelationID).
		Msg("Refresh request published - fetching document from provider")

	return nil
}
