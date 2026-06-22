package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func (uc *useCase) SyncCancellations(ctx context.Context, dto *dtos.CompareRequestDTO) (string, error) {
	dateFrom, err := time.Parse("2006-01-02", dto.DateFrom)
	if err != nil {
		return "", fmt.Errorf("invalid date_from format, expected YYYY-MM-DD: %w", err)
	}
	dateTo, err := time.Parse("2006-01-02", dto.DateTo)
	if err != nil {
		return "", fmt.Errorf("invalid date_to format, expected YYYY-MM-DD: %w", err)
	}
	if dateTo.Before(dateFrom) {
		return "", fmt.Errorf("date_to must be after date_from")
	}
	if dateTo.Sub(dateFrom) > 180*24*time.Hour {
		return "", errors.ErrCompareDateRangeTooLarge
	}

	config, err := uc.repo.GetAnyConfigByBusiness(ctx, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("business_id", dto.BusinessID).Msg("Failed to get invoicing config")
		return "", errors.ErrProviderNotConfigured
	}
	if config == nil {
		uc.log.Warn(ctx).Uint("business_id", dto.BusinessID).Msg("No invoicing config found for business")
		return "", errors.ErrProviderNotConfigured
	}

	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		uc.log.Error(ctx).Msg("No invoicing integration configured in active config")
		return "", errors.ErrProviderNotConfigured
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Warn(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to resolve provider, defaulting to softpymes")
		provider = dtos.ProviderSoftpymes
	}

	syncConfig := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			syncConfig[k] = v
		}
	}
	syncConfig["is_testing"] = config.IsTesting
	syncConfig["base_url"] = config.BaseURL
	syncConfig["base_url_test"] = config.BaseURLTest
	syncConfig["date_from"] = dto.DateFrom
	syncConfig["date_to"] = dto.DateTo
	syncConfig["business_id"] = dto.BusinessID
	syncConfig["mode"] = "sync"

	correlationID := uuid.New().String()

	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationCompare,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config:        syncConfig,
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("provider", provider).
			Str("correlation_id", correlationID).
			Msg("Failed to publish sync cancellations request to queue")
		return "", fmt.Errorf("failed to publish sync cancellations request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Str("date_from", dto.DateFrom).
		Str("date_to", dto.DateTo).
		Str("provider", provider).
		Str("correlation_id", correlationID).
		Msg("Invoice cancellations sync request published")

	return correlationID, nil
}
