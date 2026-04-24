package dtos

type ListLotsParams struct {
	BusinessID     uint
	ProductID      string
	Status         string
	ExpiringInDays int
	Page           int
	PageSize       int
}

func (p ListLotsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListSerialsParams struct {
	BusinessID uint
	ProductID  string
	LotID      *uint
	StateID    *uint
	LocationID *uint
	Page       int
	PageSize   int
}

func (p ListSerialsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListProductUoMParams struct {
	BusinessID uint
	ProductID  string
}

type ChangeInventoryStateTxParams struct {
	LevelID       uint
	FromStateCode string
	ToStateCode   string
	Quantity      int
	Reason        string
	BusinessID    uint
	CreatedByID   *uint
}
