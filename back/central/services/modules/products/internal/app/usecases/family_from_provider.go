package usecases

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"gorm.io/datatypes"
)

type varianteDeProveedor struct {
	FamilyID   *uint
	Label      string
	Attributes datatypes.JSON
}

func (uc *UseCases) resolverVariante(ctx context.Context, dto *domain.ProductProviderUpsertDTO) varianteDeProveedor {
	vacia := varianteDeProveedor{}
	if dto == nil || len(dto.VariantAttributes) == 0 {
		return vacia
	}

	crudos, err := json.Marshal(dto.VariantAttributes)
	if err != nil {
		return vacia
	}
	canonicos, firma, err := domain.CanonicalizeVariantAttributes(datatypes.JSON(crudos))
	if err != nil || firma == "" {
		return vacia
	}

	resultado := varianteDeProveedor{
		Label:      strings.TrimSpace(dto.VariantLabel),
		Attributes: canonicos,
	}

	nombre := strings.TrimSpace(dto.FamilyName)
	if nombre == "" {
		return resultado
	}

	familia, err := uc.repo.GetProductFamilyByName(ctx, dto.BusinessID, nombre)
	if err != nil {
		return resultado
	}
	if familia == nil {
		nueva := &domain.ProductFamily{
			BusinessID: dto.BusinessID,
			Name:       nombre,
			Status:     "active",
			IsActive:   true,
		}
		if err := uc.repo.CreateProductFamily(ctx, nueva); err != nil {
			return resultado
		}
		familia = nueva
	}

	if familia.ID == 0 {
		return resultado
	}

	ocupada, err := uc.repo.VariantExistsInFamily(ctx, dto.BusinessID, familia.ID, firma, nil)
	if err != nil || ocupada {
		return resultado
	}

	id := familia.ID
	resultado.FamilyID = &id
	return resultado
}
