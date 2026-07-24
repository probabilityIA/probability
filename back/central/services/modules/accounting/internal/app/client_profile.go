package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
)

func (uc *UseCase) GetClientProfile(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
	profile, err := uc.repo.GetClientProfile(ctx, businessID)
	if err != nil {
		return nil, err
	}
	if !profile.Configured {
		return profile, nil
	}
	taxes, err := uc.repo.ListTaxes(ctx)
	if err != nil {
		return profile, nil
	}
	suggested := []uint{}
	for _, tax := range taxes {
		if !tax.IsActive {
			continue
		}
		switch tax.Code {
		case "IVA":
			if profile.AppliesIVA {
				suggested = append(suggested, tax.ID)
			}
		case "RETEFUENTE":
			if profile.AppliesReteFuente {
				suggested = append(suggested, tax.ID)
			}
		case "RETEICA":
			if profile.AppliesReteICA {
				suggested = append(suggested, tax.ID)
			}
		}
	}
	profile.SuggestedTaxIDs = suggested
	return profile, nil
}
