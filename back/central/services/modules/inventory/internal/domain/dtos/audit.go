package dtos

import "time"

type ListCycleCountPlansParams struct {
	BusinessID  uint
	WarehouseID *uint
	ActiveOnly  bool
	Page        int
	PageSize    int
}

func (p ListCycleCountPlansParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListCycleCountTasksParams struct {
	BusinessID  uint
	WarehouseID *uint
	PlanID      *uint
	Status      string
	Page        int
	PageSize    int
}

func (p ListCycleCountTasksParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListCycleCountLinesParams struct {
	BusinessID uint
	TaskID     uint
	Status     string
	Page       int
	PageSize   int
}

func (p ListCycleCountLinesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type ListDiscrepanciesParams struct {
	BusinessID uint
	TaskID     *uint
	Status     string
	Page       int
	PageSize   int
}

func (p ListDiscrepanciesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

type KardexQueryParams struct {
	BusinessID  uint
	ProductID   string
	WarehouseID uint
	From        *time.Time
	To          *time.Time
}

type ApproveDiscrepancyTxParams struct {
	BusinessID     uint
	DiscrepancyID  uint
	ReviewerID     uint
	Notes          string
	MovementTypeID uint
}
