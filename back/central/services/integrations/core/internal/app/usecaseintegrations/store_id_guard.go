package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func (uc *IntegrationUseCase) ensureStoreIDNotInUse(ctx context.Context, storeID string, integrationTypeID uint, businessID *uint, excludeID uint) error {
	if storeID == "" || integrationTypeID == 0 {
		return nil
	}

	ownerBusinessID, err := uc.repo.FindStoreIDOwner(ctx, storeID, integrationTypeID, excludeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("store_id", storeID).
			Uint("integration_type_id", integrationTypeID).
			Msg("Error al verificar si la cuenta del canal ya esta en uso")
		return err
	}

	if ownerBusinessID == nil {
		return nil
	}

	if businessID != nil && *ownerBusinessID == *businessID {
		return nil
	}

	uc.log.Warn(ctx).
		Str("store_id", storeID).
		Uint("integration_type_id", integrationTypeID).
		Uint("owner_business_id", *ownerBusinessID).
		Msg("Intento de conectar una cuenta de canal que ya pertenece a otro negocio")

	return fmt.Errorf("%w (store_id %s)", domain.ErrIntegrationStoreIDInUse, storeID)
}
