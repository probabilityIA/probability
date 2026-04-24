package dtos

type ListStockMovementTypesParams struct {
	BusinessID uint
	ActiveOnly bool
	Page       int
	PageSize   int
}

func (p ListStockMovementTypesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}
