package repository

import (
	"context"
	"sort"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type familyRow struct {
	ID             uint
	Name           string
	ParentFamilyID *uint
}

func (r *Repository) GetCatalogFilters(ctx context.Context, businessID uint) (entities.StorefrontFilters, error) {
	var categories []string
	if err := r.db.Conn(ctx).Model(&models.Product{}).
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL AND category <> ''", businessID).
		Distinct("category").
		Order("category ASC").
		Pluck("category", &categories).Error; err != nil {
		return entities.StorefrontFilters{}, err
	}

	var families []familyRow
	if err := r.db.Conn(ctx).
		Table("product_families").
		Select("DISTINCT product_families.id, product_families.name, product_families.parent_family_id").
		Joins("JOIN products ON products.family_id = product_families.id AND products.deleted_at IS NULL AND products.is_active = true").
		Where("product_families.business_id = ? AND product_families.is_active = true AND product_families.deleted_at IS NULL", businessID).
		Scan(&families).Error; err != nil {
		return entities.StorefrontFilters{}, err
	}

	familyByID := make(map[uint]familyRow, len(families))
	var missingParentIDs []uint
	for _, f := range families {
		familyByID[f.ID] = f
	}
	for _, f := range families {
		if f.ParentFamilyID != nil {
			if _, ok := familyByID[*f.ParentFamilyID]; !ok {
				missingParentIDs = append(missingParentIDs, *f.ParentFamilyID)
			}
		}
	}

	if len(missingParentIDs) > 0 {
		var parents []familyRow
		if err := r.db.Conn(ctx).
			Table("product_families").
			Select("id, name, parent_family_id").
			Where("id IN ? AND business_id = ? AND is_active = true AND deleted_at IS NULL", missingParentIDs, businessID).
			Scan(&parents).Error; err != nil {
			return entities.StorefrontFilters{}, err
		}
		for _, p := range parents {
			if _, ok := familyByID[p.ID]; !ok {
				familyByID[p.ID] = p
				families = append(families, p)
			}
		}
	}

	sort.Slice(families, func(i, j int) bool {
		iHasParent := families[i].ParentFamilyID != nil
		jHasParent := families[j].ParentFamilyID != nil
		if iHasParent != jHasParent {
			return !iHasParent
		}
		return families[i].Name < families[j].Name
	})

	result := entities.StorefrontFilters{Categories: categories}
	for _, f := range families {
		result.Families = append(result.Families, entities.StorefrontFamily{
			ID:             f.ID,
			Name:           f.Name,
			ParentFamilyID: f.ParentFamilyID,
		})
	}
	return result, nil
}
