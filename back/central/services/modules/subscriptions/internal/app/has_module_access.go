package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) HasModuleAccess(ctx context.Context, businessID uint, moduleCode string) (bool, error) {
	overrides, err := uc.repo.ListOverridesByBusiness(ctx, businessID)
	if err != nil {
		return false, err
	}
	for _, o := range overrides {
		if o.ModuleCode == moduleCode {
			return true, nil
		}
	}

	if moduleregistry.IsRestrictedByDefault(moduleCode) {
		return false, nil
	}

	subTypeID, err := uc.repo.GetBusinessCurrentSubscriptionTypeID(ctx, businessID)
	if err != nil {
		return false, err
	}
	if subTypeID == nil {
		return uc.allowWhenNoPlan(ctx, businessID)
	}

	subType, err := uc.repo.GetSubscriptionType(ctx, *subTypeID)
	if err != nil {
		return false, err
	}
	if subType == nil {
		return uc.allowWhenNoPlan(ctx, businessID)
	}

	for _, code := range subType.ModuleCodes {
		if code == moduleCode {
			return true, nil
		}
	}

	return false, nil
}

// allowWhenNoPlan resuelve el fallback para un negocio sin plan valido (nunca
// asignado, o apuntando a un plan borrado/inexistente). Si el negocio esta
// activo, se mantiene el default historico de dejar ver los modulos no
// restringidos por defecto; si esta expired/cancelled, se cierra el acceso.
func (uc *UseCase) allowWhenNoPlan(ctx context.Context, businessID uint) (bool, error) {
	meta, err := uc.repo.GetBusinessSubscriptionMeta(ctx, businessID)
	if err != nil {
		return false, err
	}
	if meta != nil && meta.Status != entities.BusinessStatusActive {
		return false, nil
	}
	return true, nil
}
