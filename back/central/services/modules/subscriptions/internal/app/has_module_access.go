package app

import (
	"context"

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
		return true, nil
	}

	subType, err := uc.repo.GetSubscriptionType(ctx, *subTypeID)
	if err != nil {
		return false, err
	}
	if subType == nil {
		return true, nil
	}

	for _, code := range subType.ModuleCodes {
		if code == moduleCode {
			return true, nil
		}
	}

	return false, nil
}
