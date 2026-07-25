package usecasebusiness

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
)

func (uc *BusinessUseCase) GetFiscalProfile(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
	if _, err := uc.repository.GetBusinessByID(ctx, businessID); err != nil {
		return nil, err
	}
	profile, err := uc.repository.GetFiscalProfile(ctx, businessID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return &domain.FiscalProfile{
			BusinessID:     businessID,
			DocumentType:   "NIT",
			PersonType:     "JURIDICA",
			IVARate:        19,
			ReteFuenteRate: 11,
			ReteICARate:    0.966,
		}, nil
	}
	return profile, nil
}

func (uc *BusinessUseCase) UpsertFiscalProfile(ctx context.Context, profile *domain.FiscalProfile) (*domain.FiscalProfile, error) {
	if _, err := uc.repository.GetBusinessByID(ctx, profile.BusinessID); err != nil {
		return nil, err
	}
	profile.DocumentType = strings.ToUpper(strings.TrimSpace(profile.DocumentType))
	if profile.DocumentType == "" {
		profile.DocumentType = "NIT"
	}
	profile.PersonType = strings.ToUpper(strings.TrimSpace(profile.PersonType))
	if profile.PersonType != "JURIDICA" && profile.PersonType != "NATURAL" {
		return nil, fmt.Errorf("person_type debe ser JURIDICA o NATURAL")
	}
	if profile.IVARate < 0 || profile.IVARate > 100 || profile.ReteFuenteRate < 0 || profile.ReteFuenteRate > 100 || profile.ReteICARate < 0 || profile.ReteICARate > 100 {
		return nil, fmt.Errorf("las tarifas deben estar entre 0 y 100")
	}
	if err := uc.repository.UpsertFiscalProfile(ctx, profile); err != nil {
		return nil, err
	}
	return uc.repository.GetFiscalProfile(ctx, profile.BusinessID)
}
