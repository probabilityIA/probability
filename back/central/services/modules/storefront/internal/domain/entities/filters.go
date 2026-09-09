package entities

type StorefrontFamily struct {
	ID             uint
	Name           string
	ParentFamilyID *uint
}

type StorefrontCatalogLayout struct {
	Columns int
	Rows    int
}

type StorefrontBanner struct {
	Enabled  bool
	ImageURL string
}

type StorefrontFilters struct {
	Categories []string
	Families   []StorefrontFamily
	Layout     StorefrontCatalogLayout
	Banner     StorefrontBanner
}
