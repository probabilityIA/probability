package request

import "time"

type CreateCountPlanDTO struct {
	BusinessID    uint
	WarehouseID   uint
	Name          string
	Strategy      string
	FrequencyDays int
	NextRunAt     *time.Time
	IsActive      bool
}

type UpdateCountPlanDTO struct {
	ID            uint
	BusinessID    uint
	WarehouseID   *uint
	Name          string
	Strategy      string
	FrequencyDays *int
	NextRunAt     *time.Time
	IsActive      *bool
}

type GenerateCountTaskDTO struct {
	BusinessID uint
	PlanID     uint
	ScopeType  string
	ScopeID    *uint
}

type StartCountTaskDTO struct {
	BusinessID uint
	TaskID     uint
	UserID     uint
}

type SubmitCountLineDTO struct {
	BusinessID uint
	LineID     uint
	CountedQty int
	UserID     *uint
}

type ApproveDiscrepancyDTO struct {
	BusinessID    uint
	DiscrepancyID uint
	ReviewerID    uint
	Notes         string
}

type RejectDiscrepancyDTO struct {
	BusinessID    uint
	DiscrepancyID uint
	ReviewerID    uint
	Reason        string
}

type KardexExportDTO struct {
	BusinessID  uint
	ProductID   string
	WarehouseID uint
	From        *time.Time
	To          *time.Time
}
