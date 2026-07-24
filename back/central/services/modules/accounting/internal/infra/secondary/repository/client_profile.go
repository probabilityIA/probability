package repository

import (
	"context"
	stderrors "errors"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"gorm.io/gorm"
)

type fiscalProfileRow struct {
	DocumentType      string
	DocumentNumber    string
	DV                string
	PersonType        string
	MunicipalityID    string
	Address           string
	Phone             string
	BillingEmail      string
	AppliesIVA        bool
	AppliesReteFuente bool
	AppliesReteICA    bool
	IVARate           float64
	ReteFuenteRate    float64
	ReteICARate       float64
}

func (r *Repository) GetClientProfile(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
	var row fiscalProfileRow
	err := r.db.Conn(ctx).Table("business_fiscal_profiles").
		Select("document_type, document_number, dv, person_type, municipality_id, address, phone, billing_email, applies_iva, applies_rete_fuente AS applies_rete_fuente, applies_rete_ica AS applies_rete_ica, iva_rate, rete_fuente_rate, rete_ica_rate").
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Limit(1).First(&row).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return &dtos.ClientProfileDTO{BusinessID: businessID, Configured: false, SuggestedTaxIDs: []uint{}}, nil
		}
		return nil, err
	}
	return &dtos.ClientProfileDTO{
		BusinessID:        businessID,
		Configured:        true,
		DocumentType:      row.DocumentType,
		DocumentNumber:    row.DocumentNumber,
		DV:                row.DV,
		PersonType:        row.PersonType,
		MunicipalityID:    row.MunicipalityID,
		Address:           row.Address,
		Phone:             row.Phone,
		BillingEmail:      row.BillingEmail,
		AppliesIVA:        row.AppliesIVA,
		IVARate:           row.IVARate,
		AppliesReteFuente: row.AppliesReteFuente,
		ReteFuenteRate:    row.ReteFuenteRate,
		AppliesReteICA:    row.AppliesReteICA,
		ReteICARate:       row.ReteICARate,
		SuggestedTaxIDs:   []uint{},
	}, nil
}
