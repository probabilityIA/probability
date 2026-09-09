package response

import "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"

type FamilyResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ParentFamilyID *uint  `json:"parent_family_id"`
}

type CatalogLayoutResponse struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
}

type CatalogBannerResponse struct {
	Enabled  bool   `json:"enabled"`
	ImageURL string `json:"image_url"`
}

type CatalogFiltersResponse struct {
	Categories []string              `json:"categories"`
	Families   []FamilyResponse      `json:"families"`
	Layout     CatalogLayoutResponse `json:"layout"`
	Banner     CatalogBannerResponse `json:"banner"`
}

func CatalogFiltersFromEntity(f entities.StorefrontFilters, imageURLBase string) CatalogFiltersResponse {
	families := make([]FamilyResponse, len(f.Families))
	for i, fam := range f.Families {
		families[i] = FamilyResponse{ID: fam.ID, Name: fam.Name, ParentFamilyID: fam.ParentFamilyID}
	}
	categories := f.Categories
	if categories == nil {
		categories = []string{}
	}
	return CatalogFiltersResponse{
		Categories: categories,
		Families:   families,
		Layout:     CatalogLayoutResponse{Columns: f.Layout.Columns, Rows: f.Layout.Rows},
		Banner:     CatalogBannerFromEntity(f.Banner, imageURLBase),
	}
}

func CatalogLayoutFromEntity(l entities.StorefrontCatalogLayout) CatalogLayoutResponse {
	return CatalogLayoutResponse{Columns: l.Columns, Rows: l.Rows}
}

func CatalogBannerFromEntity(b entities.StorefrontBanner, imageURLBase string) CatalogBannerResponse {
	return CatalogBannerResponse{Enabled: b.Enabled, ImageURL: buildFullImageURL(b.ImageURL, imageURLBase)}
}
