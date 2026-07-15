package dtos

import "time"

type ReportFilter struct {
	BusinessID uint
	StartDate  time.Time
	EndDate    time.Time
	Carrier    string
}

type OrdersFilter struct {
	BusinessID uint
	StartDate  time.Time
	EndDate    time.Time
	Carrier    string
	Collected  *bool
	HasGuide   *bool
	Search     string
	Page       int
	PageSize   int
}

type SaveCarrierConfigDTO struct {
	BusinessID         uint
	CarrierName        string
	DiscountPercentage float64
	IsActive           bool
}

type ConfirmCutDTO struct {
	BusinessID  uint
	PeriodStart time.Time
	PeriodEnd   time.Time
	OrderIDs    []string
	UserID      uint
	UserName    string
}

type SelectableOrdersFilter struct {
	BusinessID  uint
	PeriodStart time.Time
	PeriodEnd   time.Time
}
