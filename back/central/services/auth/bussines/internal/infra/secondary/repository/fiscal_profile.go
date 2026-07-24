package repository

import (
	"context"
	stderrors "errors"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func fiscalProfileToDomain(m *models.BusinessFiscalProfile) *domain.FiscalProfile {
	return &domain.FiscalProfile{
		ID:                m.ID,
		BusinessID:        m.BusinessID,
		DocumentType:      m.DocumentType,
		DocumentNumber:    m.DocumentNumber,
		DV:                m.DV,
		PersonType:        m.PersonType,
		TaxRegime:         m.TaxRegime,
		MunicipalityID:    m.MunicipalityID,
		Address:           m.Address,
		Phone:             m.Phone,
		BillingEmail:      m.BillingEmail,
		AppliesIVA:        m.AppliesIVA,
		IVARate:           m.IVARate,
		AppliesReteFuente: m.AppliesReteFuente,
		ReteFuenteRate:    m.ReteFuenteRate,
		AppliesReteICA:    m.AppliesReteICA,
		ReteICARate:       m.ReteICARate,
		Notes:             m.Notes,
	}
}

func (r *Repository) GetFiscalProfile(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
	var model models.BusinessFiscalProfile
	err := r.database.Conn(ctx).Where("business_id = ?", businessID).First(&model).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return fiscalProfileToDomain(&model), nil
}

func (r *Repository) UpsertFiscalProfile(ctx context.Context, profile *domain.FiscalProfile) error {
	updates := map[string]interface{}{
		"document_type":       profile.DocumentType,
		"document_number":     profile.DocumentNumber,
		"dv":                  profile.DV,
		"person_type":         profile.PersonType,
		"tax_regime":          profile.TaxRegime,
		"municipality_id":     profile.MunicipalityID,
		"address":             profile.Address,
		"phone":               profile.Phone,
		"billing_email":       profile.BillingEmail,
		"applies_iva":         profile.AppliesIVA,
		"iva_rate":            profile.IVARate,
		"applies_rete_fuente": profile.AppliesReteFuente,
		"rete_fuente_rate":    profile.ReteFuenteRate,
		"applies_rete_ica":    profile.AppliesReteICA,
		"rete_ica_rate":       profile.ReteICARate,
		"notes":               profile.Notes,
	}
	var existing models.BusinessFiscalProfile
	err := r.database.Conn(ctx).Where("business_id = ?", profile.BusinessID).First(&existing).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			model := &models.BusinessFiscalProfile{
				BusinessID:        profile.BusinessID,
				DocumentType:      profile.DocumentType,
				DocumentNumber:    profile.DocumentNumber,
				DV:                profile.DV,
				PersonType:        profile.PersonType,
				TaxRegime:         profile.TaxRegime,
				MunicipalityID:    profile.MunicipalityID,
				Address:           profile.Address,
				Phone:             profile.Phone,
				BillingEmail:      profile.BillingEmail,
				AppliesIVA:        profile.AppliesIVA,
				IVARate:           profile.IVARate,
				AppliesReteFuente: profile.AppliesReteFuente,
				ReteFuenteRate:    profile.ReteFuenteRate,
				AppliesReteICA:    profile.AppliesReteICA,
				ReteICARate:       profile.ReteICARate,
				Notes:             profile.Notes,
			}
			return r.database.Conn(ctx).Create(model).Error
		}
		return err
	}
	return r.database.Conn(ctx).Model(&models.BusinessFiscalProfile{}).Where("id = ?", existing.ID).Updates(updates).Error
}
