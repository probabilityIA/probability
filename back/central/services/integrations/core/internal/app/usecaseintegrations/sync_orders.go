package usecaseintegrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// SyncOrdersByIntegrationID sincroniza órdenes para una integración específica.
func (uc *IntegrationUseCase) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByIntegrationID")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	return provider.SyncOrdersByIntegrationID(ctx, integrationID)
}

// SyncOrdersByIntegrationIDWithParams sincroniza órdenes con parámetros de filtrado.
// Si la integración no soporta params, hace fallback a SyncOrdersByIntegrationID.
func (uc *IntegrationUseCase) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByIntegrationIDWithParams")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	// Intentar con parámetros; si no está soportado, fallback sin params
	err = provider.SyncOrdersByIntegrationIDWithParams(ctx, integrationID, params)
	if errors.Is(err, domain.ErrNotSupported) {
		return provider.SyncOrdersByIntegrationID(ctx, integrationID)
	}
	return err
}

// SyncOrdersByBusiness sincroniza órdenes de todas las integraciones activas de un negocio
// que soporten sincronización.
func (uc *IntegrationUseCase) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	ctx = log.WithFunctionCtx(ctx, "SyncOrdersByBusiness")

	businessIDPtr := &businessID
	filters := domain.IntegrationFilters{
		BusinessID: businessIDPtr,
		IsActive:   boolPtr(true),
	}

	integrations, _, err := uc.repo.ListIntegrations(ctx, filters)
	if err != nil {
		return fmt.Errorf("error al obtener integraciones: %w", err)
	}

	for _, integration := range integrations {
		if integration.IntegrationType == nil {
			continue
		}

		integrationTypeID := int(integration.IntegrationTypeID)
		impl, registered := uc.providerReg.Get(integrationTypeID)
		if !registered {
			continue
		}

		integrationIDStr := fmt.Sprintf("%d", integration.ID)
		err := impl.SyncOrdersByIntegrationID(ctx, integrationIDStr)
		if err != nil {
			// Skip silently — ErrNotSupported means this provider doesn't sync orders
			continue
		}
	}

	return nil
}
