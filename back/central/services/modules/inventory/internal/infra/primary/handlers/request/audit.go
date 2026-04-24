package request

import "time"

type CreateCountPlanBody struct {
	WarehouseID   uint       `json:"warehouse_id" binding:"required,min=1"`
	Name          string     `json:"name" binding:"required"`
	Strategy      string     `json:"strategy"`
	FrequencyDays int        `json:"frequency_days"`
	NextRunAt     *time.Time `json:"next_run_at"`
	IsActive      *bool      `json:"is_active"`
}

type UpdateCountPlanBody struct {
	WarehouseID   *uint      `json:"warehouse_id"`
	Name          string     `json:"name"`
	Strategy      string     `json:"strategy"`
	FrequencyDays *int       `json:"frequency_days"`
	NextRunAt     *time.Time `json:"next_run_at"`
	IsActive      *bool      `json:"is_active"`
}

type GenerateCountTaskBody struct {
	PlanID    uint   `json:"plan_id" binding:"required,min=1"`
	ScopeType string `json:"scope_type"`
	ScopeID   *uint  `json:"scope_id"`
}

type StartCountTaskBody struct {
	UserID uint `json:"user_id" binding:"required,min=1"`
}

type SubmitCountLineBody struct {
	CountedQty int `json:"counted_qty" binding:"required"`
}

type ApproveDiscrepancyBody struct {
	Notes string `json:"notes"`
}

type RejectDiscrepancyBody struct {
	Reason string `json:"reason"`
}
