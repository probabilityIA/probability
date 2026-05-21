package dtos

type ListClientGroupsParams struct {
	BusinessID uint
	Search     string
	Page       int
	PageSize   int
}

func (p ListClientGroupsParams) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

type SaveClientGroupDTO struct {
	ID          uint
	BusinessID  uint
	Name        string
	Description string
	IsActive    bool
}

type ListGroupMembersParams struct {
	BusinessID    uint
	ClientGroupID uint
	Search        string
	Page          int
	PageSize      int
}

func (p ListGroupMembersParams) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

type ListAvailableClientsParams struct {
	BusinessID  uint
	Search      string
	OnlyUngrouped bool
	Page        int
	PageSize    int
}

func (p ListAvailableClientsParams) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

type AddGroupMembersDTO struct {
	BusinessID    uint
	ClientGroupID uint
	ClientIDs     []uint
}

type CatalogPriceTarget struct {
	BusinessID    uint
	ClientGroupID *uint
	ClientID      *uint
}

type ListCatalogPricesParams struct {
	Target   CatalogPriceTarget
	Search   string
	Page     int
	PageSize int
}

func (p ListCatalogPricesParams) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

type SaveCatalogPriceItem struct {
	ProductID string
	Price     *float64
}

type SaveCatalogPricesDTO struct {
	Target CatalogPriceTarget
	Items  []SaveCatalogPriceItem
}

type EffectivePriceParams struct {
	BusinessID uint
	ProductID  string
	ClientID   uint
}
