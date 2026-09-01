package entities

import "time"

type AdminBusinessRow struct {
	ID                uint
	Name              string
	Code              string
	PlanName          string
	Status            string
	CycleStartDate    *time.Time
	CycleEndDate      *time.Time
	LastPaymentAmount *float64
	LastPaymentDate   *time.Time
	ForecastedPayment *float64
	CutoffDay         *int
}

type AdminKPIs struct {
	ActiveCount             int64
	ExpiringSoonCount       int64
	ExpiredOrSuspendedCount int64
	MRR                     float64
}

const (
	AdminStatusFilterActive       = "active"
	AdminStatusFilterExpiringSoon = "expiring_soon"
	AdminStatusFilterExpired      = "expired"
	AdminStatusFilterCancelled    = "cancelled"
	AdminStatusFilterNoPlan       = "no_plan"
)
