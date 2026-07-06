package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func (uc *useCase) RequestInventorySync(ctx context.Context, businessID, integrationID uint) (string, error) {
	if integrationID == 0 {
		return "", fmt.Errorf("integration_id is required")
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to resolve provider for inventory sync")
		return "", errors.ErrProviderNotConfigured
	}
	if provider != dtos.ProviderSiigo {
		return "", fmt.Errorf("inventory sync solo esta disponible para integraciones Siigo")
	}

	correlationID := uuid.New().String()

	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationInventorySync,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config: map[string]interface{}{
				"business_id": businessID,
			},
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Str("correlation_id", correlationID).
			Msg("Failed to publish inventory sync request to queue")
		return "", fmt.Errorf("failed to publish inventory sync request: %w", err)
	}

	uc.log.Info(ctx).
		Uint("business_id", businessID).
		Uint("integration_id", integrationID).
		Str("correlation_id", correlationID).
		Msg("Inventory sync request published")

	return correlationID, nil
}

func (uc *useCase) RequestListSiigoWarehouses(ctx context.Context, businessID, integrationID uint) (string, error) {
	if integrationID == 0 {
		return "", fmt.Errorf("integration_id is required")
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_id", integrationID).Msg("Failed to resolve provider for list siigo warehouses")
		return "", errors.ErrProviderNotConfigured
	}
	if provider != dtos.ProviderSiigo {
		return "", fmt.Errorf("el listado de bodegas solo esta disponible para integraciones Siigo")
	}

	correlationID := uuid.New().String()

	requestMessage := &dtos.InvoiceRequestMessage{
		InvoiceID: 0,
		Provider:  provider,
		Operation: dtos.OperationListSiigoWarehouses,
		InvoiceData: dtos.InvoiceData{
			IntegrationID: integrationID,
			Config: map[string]interface{}{
				"business_id": businessID,
			},
		},
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
	}

	if err := uc.invoiceRequestPub.PublishInvoiceRequest(ctx, requestMessage); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Str("correlation_id", correlationID).
			Msg("Failed to publish list siigo warehouses request")
		return "", fmt.Errorf("failed to publish list siigo warehouses request: %w", err)
	}

	return correlationID, nil
}
