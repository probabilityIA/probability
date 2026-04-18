package dtos

type ListPutawayRulesParams struct {
	BusinessID uint
	ActiveOnly bool
	Page       int
	PageSize   int
}

func (p ListPutawayRulesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListPutawaySuggestionsParams struct {
	BusinessID uint
	Status     string
	Page       int
	PageSize   int
}

func (p ListPutawaySuggestionsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListReplenishmentTasksParams struct {
	BusinessID  uint
	WarehouseID *uint
	Status      string
	AssignedTo  *uint
	Page        int
	PageSize    int
}

func (p ListReplenishmentTasksParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListCrossDockLinksParams struct {
	BusinessID      uint
	OutboundOrderID string
	Status          string
	Page            int
	PageSize        int
}

func (p ListCrossDockLinksParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListVelocityParams struct {
	BusinessID  uint
	WarehouseID uint
	Period      string
	Rank        string
	Limit       int
}
