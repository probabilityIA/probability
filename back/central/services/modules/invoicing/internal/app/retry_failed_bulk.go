package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

const (
	retryFailedMaxInvoices = 2000
	retryFailedMaxDaysBack = 35
)

func (uc *useCase) RetryFailedBulk(ctx context.Context, businessID uint) (int, error) {
	if businessID == 0 {
		return 0, fmt.Errorf("business_id is required")
	}

	invoices, _, err := uc.repo.ListInvoices(ctx, map[string]interface{}{
		"business_id": businessID,
		"status":      "failed",
		"page":        1,
		"page_size":   retryFailedMaxInvoices,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list failed invoices: %w", err)
	}

	pairs := make([]map[string]interface{}, 0, len(invoices))
	oldest := time.Now()
	for _, inv := range invoices {
		if inv.InvoiceNumber != "" {
			continue
		}
		pairs = append(pairs, map[string]interface{}{
			"invoice_id": inv.ID,
			"order_id":   inv.OrderID,
		})
		if inv.CreatedAt.Before(oldest) {
			oldest = inv.CreatedAt
		}
	}

	if len(pairs) == 0 {
		return 0, nil
	}

	config, err := uc.repo.GetAnyConfigByBusiness(ctx, businessID)
	if err != nil || config == nil {
		return 0, errors.ErrProviderNotConfigured
	}

	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		return 0, errors.ErrProviderNotConfigured
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		provider = dtos.ProviderSoftpymes
	}
	if provider != dtos.ProviderSoftpymes {
		return 0, fmt.Errorf("reintento masivo de fallidas solo disponible para Softpymes por ahora (proveedor actual: %s)", provider)
	}

	minDate := time.Now().AddDate(0, 0, -retryFailedMaxDaysBack)
	if oldest.Before(minDate) {
		oldest = minDate
	}

	reconcileConfig := make(map[string]interface{})
	if config.InvoiceConfig != nil {
		for k, v := range config.InvoiceConfig {
			reconcileConfig[k] = v
		}
	}
	reconcileConfig["reconcile_invoices"] = pairs
	reconcileConfig["date_from"] = oldest.Format("2006-01-02")

	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationReconcileFailed,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config:        reconcileConfig,
		},
		CorrelationID: uuid.New().String(),
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		return 0, fmt.Errorf("failed to publish reconcile request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", businessID).
		Int("invoices", len(pairs)).
		Str("date_from", reconcileConfig["date_from"].(string)).
		Msg("Bulk retry of failed invoices requested (reconcile first)")

	return len(pairs), nil
}
