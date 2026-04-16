package dtos

// CatalogFilters for public catalog listing
type CatalogFilters struct {
	Search   string
	Category string
	Page     int
	PageSize int
}

func (f CatalogFilters) Offset() int {
	if f.Page < 1 {
		f.Page = 1
	}
	return (f.Page - 1) * f.PageSize
}

func (f *CatalogFilters) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 12
	}
}
