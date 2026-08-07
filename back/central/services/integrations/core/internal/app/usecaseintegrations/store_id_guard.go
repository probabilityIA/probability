package usecaseintegrations

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func (uc *IntegrationUseCase) ensureStoreIDNotInUse(ctx context.Context, storeID string, integrationTypeID uint, businessID *uint, excludeID uint) error {
	if storeID == "" || integrationTypeID == 0 {
		return nil
	}

	owner, err := uc.repo.FindStoreIDOwner(ctx, storeID, integrationTypeID, excludeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("store_id", storeID).
			Uint("integration_type_id", integrationTypeID).
			Msg("Error al verificar si la cuenta del canal ya esta en uso")
		return err
	}

	if owner == nil {
		return nil
	}

	sameBusiness := businessID != nil && owner.BusinessID != nil && *owner.BusinessID == *businessID

	event := uc.log.Warn(ctx).
		Str("store_id", storeID).
		Uint("integration_type_id", integrationTypeID).
		Uint("existing_integration_id", owner.IntegrationID).
		Bool("same_business", sameBusiness)
	if owner.BusinessID != nil {
		event = event.Uint("owner_business_id", *owner.BusinessID)
	}
	event.Msg("Intento de conectar una cuenta de canal que ya esta en uso")

	if sameBusiness {
		return domain.ErrIntegrationStoreIDSameBusiness
	}
	return domain.ErrIntegrationStoreIDInUse
}

func isUniqueStoreIDViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "ux_integrations_store_type")
}
