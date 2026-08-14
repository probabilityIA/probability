package usecaseproduct

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCaseProduct) MatchFamilyImportRows(ctx context.Context, businessID uint, rows []domain.DimensionRow) (*domain.FamilyImportPreview, error) {
	preview := &domain.FamilyImportPreview{
		Matched:  []domain.FamilyImportMatch{},
		NotFound: []string{},
	}

	for _, row := range rows {
		product, err := uc.repo.GetProductBySKU(ctx, businessID, row.SKU)
		if err != nil {
			preview.NotFound = append(preview.NotFound, row.SKU)
			continue
		}

		match := domain.FamilyImportMatch{
			ProductID:       product.ID,
			SKU:             product.SKU,
			Name:            product.Name,
			CurrentFamilyID: product.FamilyID,
			Weight:          row.Weight,
			Length:          row.Length,
			Width:           row.Width,
			Height:          row.Height,
		}
		if product.Family != nil {
			match.CurrentFamilyName = product.Family.Name
		}

		preview.Matched = append(preview.Matched, match)
	}

	return preview, nil
}
